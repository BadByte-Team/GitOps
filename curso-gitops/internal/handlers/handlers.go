package handlers

import (
	"curso-gitops/internal/auth"
	"curso-gitops/internal/models"
	"curso-gitops/internal/repository"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// jsonResponse sets Content-Type and encodes the response as JSON.
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// jsonError returns a JSON error response.
func jsonError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// decodeBody decodes JSON body with a size limit.
func decodeBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB limit
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		jsonError(w, http.StatusBadRequest, "Cuerpo de solicitud inválido")
		return false
	}
	return true
}

func Login(w http.ResponseWriter, r *http.Request) {
	var c models.Credentials
	if !decodeBody(w, r, &c) {
		return
	}
	c.Username = strings.TrimSpace(c.Username)
	if c.Username == "" || c.Password == "" {
		jsonError(w, http.StatusBadRequest, "Usuario y contraseña son requeridos")
		return
	}

	role, err := repository.GetUserRole(c.Username, c.Password)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, "Usuario o contraseña incorrectos")
		return
	}

	token, err := auth.GenerateJWT(c.Username, role)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Error al generar token")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"token": token, "role": role})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var c models.Credentials
	if !decodeBody(w, r, &c) {
		return
	}
	c.Username = strings.TrimSpace(c.Username)

	if len(c.Username) < 3 {
		jsonError(w, http.StatusBadRequest, "El usuario debe tener al menos 3 caracteres")
		return
	}
	if len(c.Password) < 6 {
		jsonError(w, http.StatusBadRequest, "La contraseña debe tener al menos 6 caracteres")
		return
	}
	if len(c.Username) > 50 {
		jsonError(w, http.StatusBadRequest, "El usuario no puede exceder 50 caracteres")
		return
	}

	if err := repository.CreateUser(c.Username, c.Password); err != nil {
		jsonError(w, http.StatusConflict, "El usuario ya existe o los datos son inválidos")
		return
	}
	jsonResponse(w, http.StatusCreated, map[string]string{"message": "Usuario creado exitosamente"})
}

func GetModules(w http.ResponseWriter, r *http.Request) {
	isAdmin := r.Context().Value(auth.RoleKey) == "admin"
	modules, err := repository.GetModules(isAdmin)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Error al obtener módulos")
		return
	}
	jsonResponse(w, http.StatusOK, modules)
}

func CreateModule(w http.ResponseWriter, r *http.Request) {
	var m models.Module
	if !decodeBody(w, r, &m) {
		return
	}
	m.Title = strings.TrimSpace(m.Title)
	if m.Title == "" {
		jsonError(w, http.StatusBadRequest, "El título es requerido")
		return
	}
	if err := repository.AddModule(m.Title); err != nil {
		jsonError(w, http.StatusInternalServerError, "Error al crear módulo")
		return
	}
	jsonResponse(w, http.StatusCreated, map[string]string{"message": "Módulo creado"})
}

func DeleteModule(w http.ResponseWriter, r *http.Request) {
	if err := repository.DeleteModule(chi.URLParam(r, "id")); err != nil {
		jsonError(w, http.StatusInternalServerError, "Error al eliminar módulo")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Módulo eliminado"})
}

func ToggleModule(w http.ResponseWriter, r *http.Request) {
	if err := repository.ToggleModule(chi.URLParam(r, "id")); err != nil {
		jsonError(w, http.StatusInternalServerError, "Error al actualizar módulo")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Visibilidad actualizada"})
}

func UpdateModule(w http.ResponseWriter, r *http.Request) {
	var m models.Module
	if !decodeBody(w, r, &m) {
		return
	}
	m.Title = strings.TrimSpace(m.Title)
	if m.Title == "" {
		jsonError(w, http.StatusBadRequest, "El título es requerido")
		return
	}
	if err := repository.UpdateModule(chi.URLParam(r, "id"), m.Title); err != nil {
		jsonError(w, http.StatusInternalServerError, "Error al actualizar módulo")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Módulo actualizado"})
}

func CreateEpisode(w http.ResponseWriter, r *http.Request) {
	var ep models.Episode
	if !decodeBody(w, r, &ep) {
		return
	}
	ep.Title = strings.TrimSpace(ep.Title)
	ep.VideoURL = strings.TrimSpace(ep.VideoURL)
	if ep.Title == "" || ep.VideoURL == "" || ep.ModuleID == 0 {
		jsonError(w, http.StatusBadRequest, "Título, URL y ID del módulo son requeridos")
		return
	}
	if err := repository.AddEpisode(ep.ModuleID, ep.Title, ep.VideoURL); err != nil {
		jsonError(w, http.StatusInternalServerError, "Error al crear episodio")
		return
	}
	jsonResponse(w, http.StatusCreated, map[string]string{"message": "Episodio creado"})
}

func DeleteEpisode(w http.ResponseWriter, r *http.Request) {
	if err := repository.DeleteEpisode(chi.URLParam(r, "id")); err != nil {
		jsonError(w, http.StatusInternalServerError, "Error al eliminar episodio")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Episodio eliminado"})
}

func UpdateEpisode(w http.ResponseWriter, r *http.Request) {
	var ep models.Episode
	if !decodeBody(w, r, &ep) {
		return
	}
	ep.Title = strings.TrimSpace(ep.Title)
	ep.VideoURL = strings.TrimSpace(ep.VideoURL)
	if ep.Title == "" || ep.VideoURL == "" {
		jsonError(w, http.StatusBadRequest, "Título y URL son requeridos")
		return
	}
	if err := repository.UpdateEpisode(chi.URLParam(r, "id"), ep.Title, ep.VideoURL); err != nil {
		jsonError(w, http.StatusInternalServerError, "Error al actualizar episodio")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Episodio actualizado"})
}
