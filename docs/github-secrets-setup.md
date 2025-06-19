# GitHub Secrets Setup para CI/CD Pipeline

Para que la pipeline de GitHub Actions funcione correctamente, necesitas configurar los siguientes secrets en tu repositorio.

## Cómo agregar secrets

1. Ve a tu repositorio en GitHub
2. Navega a **Settings > Secrets and variables > Actions**
3. Haz clic en **New repository secret**
4. Agrega cada uno de los siguientes secrets:

## Secrets requeridos

### 1. Docker Hub
```bash
DOCKER_USERNAME=tu_username_dockerhub
DOCKER_PASSWORD=tu_password_o_token_dockerhub
```

### 2. EC2 Access
```bash
EC2_PUBLIC_IP=3.210.158.123
EC2_SSH_PRIVATE_KEY=contenido_completo_del_archivo_pem
```

Para el `EC2_SSH_PRIVATE_KEY`, necesitas copiar todo el contenido del archivo `.pem`:

```bash
# En tu máquina local:
cat terraform/stock-recommender-key.pem
```

Copia todo el output (incluyendo las líneas `-----BEGIN RSA PRIVATE KEY-----` y `-----END RSA PRIVATE KEY-----`) y pégalo como valor del secret.

## Configuración de Docker Hub

### Opción 1: Usar tu cuenta existente
- Username: tu username de Docker Hub
- Password: tu password de Docker Hub

### Opción 2: Crear token de acceso (recomendado)
1. Ve a Docker Hub > Account Settings > Security
2. Crea un nuevo Access Token
3. Usa tu username y el token como password

## Verificar configuración

Una vez configurados los secrets, puedes:

1. **Hacer push a main** para activar la pipeline automáticamente
2. **Ejecutar manualmente** desde GitHub Actions tab

## Estructura de la pipeline

La pipeline hace:

1. **Build**: 
   - Construye imágenes Docker del frontend y backend
   - Las sube a Docker Hub con tags únicos

2. **Deploy**:
   - Se conecta a EC2 via SSH
   - Usa Ansible para deployment
   - Actualiza containers con nuevas imágenes

3. **Verify**:
   - Testa endpoints de salud
   - Verifica que la aplicación responda

## Troubleshooting

### Error de SSH
- Verifica que `EC2_SSH_PRIVATE_KEY` tenga todo el contenido del .pem
- Asegúrate que la IP en `EC2_PUBLIC_IP` sea correcta

### Error de Docker Hub
- Verifica username/password
- Asegúrate que los repositorios existan en Docker Hub

### Error de Ansible
- Verifica que Ansible inventory tenga la IP correcta
- Revisa que los puertos estén abiertos en Security Groups 