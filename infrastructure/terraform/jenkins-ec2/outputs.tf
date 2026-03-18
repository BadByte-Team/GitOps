output "jenkins_public_ip" {
  description = "IP publica del servidor Jenkins"
  value       = aws_instance.jenkins.public_ip
}

output "jenkins_instance_id" {
  description = "ID de la instancia EC2"
  value       = aws_instance.jenkins.id
}

output "jenkins_url" {
  description = "URL de Jenkins"
  value       = "http://${aws_instance.jenkins.public_ip}:8080"
}

output "sonarqube_url" {
  description = "URL de SonarQube"
  value       = "http://${aws_instance.jenkins.public_ip}:9000"
}
