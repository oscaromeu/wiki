---
id: configuracion
title: Configuracion
sidebar_label: Configuration
sidebar_position: 3
---

## Variables de entorno

Para añadir alguna configuración específica a la hora de lanzar un
contenedor, se usan variables de entorno  del contenedor
cuyos valores se especifican al crear el contenedor para realizar una configuración concreta del mismo.

Por ejemplo, si estudiamos la documentación de la imagen `mariadb` en
[Docker Hub](https://hub.docker.com/_/mariadb), podemos comprobar que
podemos definir un conjunto de variables de entorno como
`MYSQL_ROOT_PASSWORD`, `MYSQL_DATABASE`, `MYSQL_USER`,
`MYSQL_PASSWORD`, etc., que nos permitirán configurar de alguna forma
determinada nuestro servidor de base de datos (indicando la contraseña
del usuario root, creando una determinada base de datos o creando un
usuario con una contraseña por ejemplo.

De la misma manera, al especificar los contenedores que contendrán los
Pods que se van a crear desde un Deployment, también se pondrán
inicializar las variables de entorno necesarias.

### Configuración de aplicaciones usando variables de entorno

Vamos a hacer un despliegue de un servidor de base de datos
mariadb. Si volvemos a estudiar la documentación de esta imagen en
[Docker Hub](https://hub.docker.com/_/mariadb) comprobamos que
obligatoriamente tenemos que indicar la contraseña del usuario root
inicializando la variable de entorno `MYSQL_ROOT_PASSWORD`. El fichero
de despliegue que vamos a usar es
[`mariadb-deployment-env.yaml`](files/mariadb-deployment-env.yaml), y
vemos el fragmento del fichero donde se define el contenedor:

```yaml
...
    spec:
      containers:
        - name: contenedor-mariadb
          image: mariadb
          ports:
            - containerPort: 3306
              name: db-port
          env:
            - name: MYSQL_ROOT_PASSWORD
              value: my-password
```

En el apartado `containers` hemos incluido la sección `env` donde
vamos indicando, como una lista, el nombre de la variable (`name`) y
su valor (`value`). En este caso hemos indicado la contraseña
`my-password`.

Vamos a comprobar si realmente se ha creado el servidor de base de
datos con esa contraseña del root:

    kubectl apply -f mariadb-deploymen-env.yaml

    kubectl get all
    ...
    NAME                                 READY   UP-TO-DATE   AVAILABLE   AGE
    deployment.apps/mariadb-deployment   1/1     1            1           5s

    kubectl exec -it deployment.apps/mariadb-deployment -- mysql -u root -p
    Enter password:
    ...
    MariaDB [(none)]>


## Configuración centralizada

En el apartado anterior hemos estudiado como podemos definir las
variables de entorno de los contenedores que vamos a desplegar. Sin
embargo, la solución que presentamos puede tener alguna limitación:

* Los valores de las variables de entorno están escritos directamente
  en el fichero yaml. Estos ficheros yaml suelen estar en repositorios
  git y lógicamente no es el sitio más adecuado para ubicarlos.
* Por otro lado, escribiendo los valores de las variables de entorno
  directamente en los ficheros, hacemos que estos ficheros no sean
  reutilizables en otros despliegues y que el procedimiento de cambiar
  las variables sea tedioso y propenso a errores, porque hay que
  hacerlo en varios sitios.

Para solucionar estas limitaciones, podemos usar dos nuevos recursos de
Kubernetes: `ConfigMap` y `Secrets`

Estos dos objetos permiten configurar las aplicaciones de manera centralizada y permiten además desacoplar sus ciclos de vida. La ventaja de separar las configuraciones de los sistemas se ven claramente cuando la empresa cuenta con diferentes entornos de despliegues, por ejemplo, cualificación, producción y laboratorio. Para cada entorno tiene que funcionar la misma imagen de Docker, pero cambiando los datos de las configuraciones. 

En esencia, los objetos `ConfigMap` o `Secret` contienen pares de clave valor. Estos pares de clave valor puede ser inyectados en el contenedor como variables de entorno o pueden ser montados como volumen. 


### Configuración de aplicaciones usando ConfigMaps

#### Creando un ConfigMap

Los [ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/) los podemos crear de manera imperativa con un simple comando: `kubectl create configmap`. Como parte del comando tenemos que indicar en las opciones cuál es el origen de los datos. Kubernetes nos brinda cuatro opciones:

+ Valores literales en texto plano en formato clave-valor.
+ Un fichero que contiene datos en formato clave-valor las cuáles se espera que sean usadas como variables de entorno.
+ Un fichero con contenido arbitrario
+ Un directorio con uno o varios ficheros.

El siguiente comando muestra todas las opciones en acción:

_Valores literales_

  ```
  $ kubectl create configmap db-config --from-literal=db=staging
  ```
  
_Fichero con variables de entorno_

  ```
  $ kubectl create configmap db-config --from-env-file=config.env
  ```

_Fichero_

  ```
  $ kubectl create configmap db-config --from-file=config.txt
  ```

_Directorio que contiene ficheros_

  ```
  $ kubectl create configmap db-config --from-file=app-config
  ```

De manera declarativa el fitxero YAML tiene este aspecto

```
apiVersion: v1
kind: ConfigMap
metadata: 
  name: backend-config
data:
  database_url: jdbc:postgresql://localhost/test
  user: fred
```

### Consumir el ConfigMap como variable de entorno

Una vez el ConfigMap ha sido creado puede ser consumido por uno o varios Pods en el mismo namespace. Vamos a ver cómo podemos inyectar los pares clave-valor del configmap como variables de entorno. 

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: configured-pod
spec:
  containers:
  - image: nginx:1.19.0
    name: app
    envFrom:
    - configMapRef:
        name: backend-config
```

El pod del ejemplo anterior usa el atributo `envFrom.configMapRef` para indicar el configMap que se desea consumir y inyectar los valores como variables de entorno.

Observar que el atributo `envFrom` no formatea de manera automática las variables de entorno siguiendo ningún estandar sino que simplemente se usa el valor tal cuál está. Después de haber creado el Pod podemos inspeccionar las variables de entorno inyectadas ejecutando de manera remota el comando UNIX `env`: 

```
$ kubectl exec configured-pod -- env
...
database_url: jdbc:postgresql://localhost/test
user=fred
```

Si necesitamos que los pares clave-valor tengan algún tipo de formato podemos redefinir las claves inyectadas con el atributo `valueFrom`. El siguiente ejemplo redefine la clave `database_url_` en `DATABASE_URL` y la clave `user` en `USERNAME`. 

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: configured-pod
spec:
  containers:
  - image: nginx:1.19.0
    name: app
    env:
    - name: DATABASE_URL
      valueFrom:
        configMapKeyRef:
          name: backend-config
          key: database_url
    - name: USERNAME
      valueFrom:
        configMapKeyRef:
          name: backend-config
          key: user
```

Si inspeccionamos de nuevo las variables de entorno en el pod veremos que tienen el tipico aspecto de variables de entorno

```
$ kubectl exec configured-pod -- env
...
DATABASE_URL=jdb:postgresql://localhost/test
USERNAME=fred
...
```

### Montando un ConfigMap como Volumen

En múltiples ocasiones las configuraciones de las aplicaciones se establecen en ficheros y no como variables de entorno. Por ejemplo,

+ `/etc/nginx/nginx.conf`: configuraciones para el servidor nginx
+ `/etc/mysql/my.cnf`: configuraciones para el servidor de base de datos MySQL
+ `/usr/share/postgresql/postgresql.conf`: configuraciones para el servidor de base de datos PostgreSQL

Ante esta situación será necesario establecer el contenido del fichero dentro del objeto ConfigMap y luego incluirlo en el Pod como un volumen.

```
apiVersion: v1
kind: Pod
metadata:
  name: configured-pod
spec:
  containers:
  - image: nginx:1.19.0
    name: app
    volumeMounts:
    - name: config-volume
      mountPath: /etc/config
  volumes:
  - name: config-volume
    configMap:
      name: backend-config
```

El atributo `volumes` especifica el volumen a usar. Como podemos ver hace referencia al configMap. El nombre del volumen es necesario para añadir el punto de montaje utilizando el atributo `volumeMounts` el cuál apunta al directorio `/etc/config`

```
$ kubectl exec -it configured-pod -- sh
# ls -1 /etc/config
database_url
user
# cat /etc/config/database_url
jdbc:postgresql://localhost/test
# cat /etc/config/user
fred
```

## Creando secretos

La estructura del objeto Secret es muy similar a la del objeto ConfigMap. La principal diferencia es que la información tiene que estar codificada en base64. El comando Create nos va a facilitar el trabajo codificando la información antes de crear el objeto en el cluster. A través de este comando podemos crear tres tipos de objeto Secret:

|Tipo|Descripción|
|----|-----------|
|generic| Crear un secreto a partir de un archivo local, directorio o valor literal|
|tls| Crar un secreto TLS|
|docker-registry|Cree un secreto para usar con un registro de Docker. |

En la mayoria de casos vamos a trabajar con objetos Secret de tipo Generico los cuáles tienen las mismas opciones de configuración que el objeto de tipo ConfigMap. 

+ Valores literales en texto plano en formato clave-valor.
+ Un fichero que contiene datos en formato clave-valor las cuáles se espera que sean usadas como variables de entorno.
+ Un fichero con contenido arbitrario
+ Un directorio con uno o varios ficheros.

Vamos a ver algunos comandos para crear de manera imperativa secretos

_Valores literales_

  ```
  $ kubectl create secret generic db-creds --from-literal=pwd=s3cre!
  ```
  
_Fichero con variables de entorno_

  ```
  $ kubectl create secret generic db-creds --from-env-file=secret.env
  ```

_Fichero SSH_

  ```
  $ kubectl create secret generic ssh-key --from-file=id_rsa=~/.ssh/id_rsa
  ```


Si optamos por hacerlo de manera declarativa tendremos de codificar en base64 los datos independientemente de si es un fichero o una cadena de texto.

+ Para cadenas de texto: `echo -n mipassword | base64 -w0`
+ Para ficheros: `cat <fichero>|base64 -w0`

El `-w0` siempre para que todo nos venga en una linea. Y para cadenas de texto siempre metemos `echo -n` ya que si no va a poner un salto de linea (`\n`) en la password. Por ejemplo, 

```
$ echo -n 's3cre!' | base64 -w0
czNjcmUh
```

este dato ya lo podemos usar en nuestro secreto

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-creds
type: Opaque
data:
  pwd: czNjcmUh
```

consultar la documentación de kubernetes para más [información](https://kubernetes.io/es/docs/concepts/configuration/secret/).

### Consumir el secreto como variables de entorno.

Consumir los datos del secreto como variables de entorno funciona de la misma manera que para un configMap. Solamente hay una diferencia: en lugar de usar el atributo `envFrom.configMapRef` usaremos `envFrom.secretRef` tal como se muestra en el siguiente ejemplo

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: configured-pod
spec:
  containers:
  - image: nginx:1.19.0
    name: app
    envFrom:
    - secretRef:
        name: db-creds
```

```
$ kubectl exec configured-pod -- env | grep -i pwd
pwd=s3cre!
```

### Montar el secreto como volumen

Montar el secreto como volumen es muy similar a como lo hicimos para el configMap la mayor diferencia que la referencia se hace a un secreto a través del atributo `secret.secretName`

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: configured-pod
spec:
  containers:
  - image: nginx:1.19.0
    name: app
    volumeMounts:
    - name: secret-volume
      mountPath: /var/app
      readOnly: true
  volumes:
  - name: secret-volume
    secret:
      secretName: db-creds
```

```
$ kubectl get pods
NAME             READY   STATUS    RESTARTS   AGE
configured-pod   1/1     Running   0          3s

$ kubectl exec -it configured-pod -- sh
# ls -1 /var/app
pwd
# cat /var/app/pwd
s3cre!# 
```
