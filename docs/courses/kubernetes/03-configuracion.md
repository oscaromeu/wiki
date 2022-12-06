---
id: configuracion
title: Configuracion
sidebar_label: Configuration
sidebar_position: 4
draft: true
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

## Estrategias de gestión de secretos


## Configurar contextos de seguridad

## RBAC

Los clusters de Kubernetes son utilizados por múltiples grupos de personas dentro de la empresa. Una primera aproximación de esta afirmación pudiera ser que el equipo de operaciones se encargase de mantener el buen funcionamiento del cluster, mientras que el equipo de desarrollo despliega y actualiza continuamente las aplicaciones.

A esta primera división lógica se le puede aumentar la complejidad si se entiende que pueden existir múltiples equipos de operaciones, donde cada uno está especializado en un área de la infraestructura p.ej: almacenamiento, monitorización, seguridad. De igual forma los equipos de desarrollo pudieran estar divididos por clientes, donde cada grupo gestiona solamente los servicios para los que han sido contratado.

Si su empresa cuenta con múltiples grupos como los mencionados anteriormente puede ser un buen síntoma, pudiera significar un alto grado de evolución y madurez, pero para el cluster tiene una connotación distinta. Piense qué implicaciones tendría si: 

+ Un desarrollador elimina accidentalmente el servicio de monitorización.
+ Un integrante del equipo de monitorización elimina accidentalmente una regla de seguridad.
+ Un integrante del cliente A obtiene las credenciales de base de datos del cliente B

A estas y muchas otras situaciones puede estar expuesto el cluster si no se controlan correctamente los permisos con los que acceden las personas de los diferentes grupos. Cada equipo de trabajo debe tener bien definidas sus responsabilidades y por lo tanto, sus permisos en el cluster deben ser acotados para que realicen solamente sus tareas y no otras.
Para dar respuesta a esta situación Kubernetes tiene implementado un sistema de permisos basado en roles (Role Base Access Control, RBAC). Este sistema es el recomendado por la plataforma para gestionar los permisos de los usuarios y servicios que deseen acceder a los recursos del cluster. Usuarios y Servicios son los elementos que pueden desencadenar el proceso de control de acceso al cluster. Los usuarios van a corresponder con todas las peticiones que se realizan desde el exterior del cluster, por ejemplo desde Kubectl o cualquier otro cliente que se comunique con Kubernetes. Por otro lado, los servicios serán todas las peticiones que se realizan desde el interior, que en este caso son los procesos que se están ejecutando dentro de los Pods. 

El control de acceso basado en roles RBAC se lleva a cabo en el API de Kubernetes. Cuando
las peticiones llegan a este componente tendrán que pasar por múltiples fases para saber si serán ejecutadas en el cluster o no. En la documentación oficial se describen cuatro fases para realizar el proceso de control, pero en el capítulo serán analizadas las dos más importantes: Autenticación y Autorización.

La fase de Autenticación será la responsable de identificar si el usuario o servicio es válido para el cluster, mientras que la fase de Autorización evaluará si es posible realizar la acción solicitada. 

### Caso práctico
Ambas fases serán llevadas a la práctica en las próximas secciones. La situación a resolver será la siguiente:

+ El usuario foo pertenece al equipo developer y necesita poder listar los Pods.
+ Una app ha sido desplegada en el cluster y necesita los permisos para listar los Pods.

Como podrá darse cuenta, en ambos casos (usuario foo y servicio app) necesitarán tener acceso para listar los Pods en el cluster.

### Gestión de usuarios externos

Lo más importante de la fase Autenticación será entender que los usuarios externos a Kubernetes no van a corresponder con el modelo clásico de usuario y contraseña. Como tampoco va a encontrar una opción dentro de los comandos Kubectl para crear un usuario. La razón de este extraño comportamiento se debe a que no son propiamente usuarios, sino que son certificados. Kubernetes utiliza los certificados X.509 en el proceso de autenticación. Estos certificados deberán estar firmados por una Autoridad Certificadora (CA) para el cluster en el momento de ser presentados. Cuando llega la solicitud al API de Kubernetes se comprueba que el certificado sea válido, y si está correcto se pasa a utilizar el nombre (CN) del parámetro -subj como el nombre de usuario y las organizaciones (O) como grupos a los que pertenece.

:::info

La CA utilizada será el certificado generado por la instalación de Kubernetes. El resultado de este paso será la obtención de un certificado firmado y válido por un año.

:::

```
docker container cp book-control-plane:/etc/kubernetes/pki/ca.crt rbac/ca.crt
docker container cp book-control-plane:/etc/kubernetes/pki/ca.key rbac/ca.key
```

### KUBECONFIG

El fichero KUBECONFIG es el vínculo de conexión entre un cliente y el cluster. Este fichero va a contener la información necesaria para que se realice la comunicación entre ambos puntos.

KUBECONFIG ha sido definido por el tipo de objeto Config, que al igual que el resto de objetos en Kubernetes puede ser accedido y configurado a través de Kubectl. Consulte la configuración que está utilizando en este momento a través del siguiente comando:

```
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: DATA+OMITTED
    server: https://127.0.0.1:34591
  name: kind-kind
contexts:
- context:
    cluster: kind-kind
    user: kind-kind
  name: kind-kind
current-context: kind-kind
kind: Config
preferences: {}
users:
- name: kind-kind
  user:
    client-certificate-data: REDACTED
    client-key-data: REDACTED
```

Observe que existen tres grupos de información en este fichero: clusters, usuarios y contextos. Los clusters contienen la información de los servidores donde se desea acceder, mientras que los usuarios contienen los certificados utilizados para la autenticación en Kubernetes. El tercer grupo son los contextos, estos últimos serán el punto donde se unen los dos grupos anteriores. 

Un contexto va a corresponder con la unión de un usuario y un cluster, a esta unión se le establece un nombre para que pueda ser identificada. Con esta manera de asociar usuarios y clusters se podrán realizar todas las combinaciones necesarias en aras de establecer la correcta autenticación a un servidor.

Cuando se crea el cluster de Kubernetes se genera automáticamente el primer fichero de KUBECONFIG. Este fichero va a contener la información para acceder al cluster como super usuario. Luego deberán ser creados y configurados el resto de ficheros en dependencia de las necesidades de acceso de los grupos de trabajo.

Para el cluster de ejemplo que se ha iniciado con Kind se tiene el primer fichero creado en la dirección `∼/.kube/config`, y a continuación se creará un segundo fichero para que el usuario foo pueda acceder al cluster con sus credenciales. Este sería el mismo procedimiento a seguir si usted fuese el administrador del cluster y un usuario le solicita acceso para listar los Pods. Como resultado de su configuración deberá enviarle el nuevo fichero generado a su colega. 

1. Adicionar la información del cluster. Se agrega la dirección del servidor junto con su certificado. El parámetro --kubeconfig le permitirá definir el nombre del fichero KUBECONFIG. Si este fichero no existe en el directorio
entonces será creado. 

```
kubectl --kubeconfig=rbac/foo-config \
config set-cluster kubernetes \
--server https://127.0.0.1:6443 \
--certificate-authority=rbac/ca.crt \
--embed-certs=true
Cluster "kubernetes" set.
```

El parámetro --embed-certs va a codificar en Base64 el certificado utilizado en el campo
certificate-authority.

2. Adicionar el usuario con sus certificados En este paso se incluyen los certificados creados en la sección anterior.

```
kubectl --kubeconfig=rbac/foo-config \
config set-credentials foo \
--client-certificate=rbac/foo.crt \
--client-key=rbac/foo.key \
--embed-certs=true
User "foo" set.
```

3. Crear el contexto. Se establece la relación entre el usuario foo y el servidor. También se ha definido el
Namespace a utilizar en caso de no ser establecido como parámetro en las líneas de comandos.

```
kubectl --kubeconfig=rbac/foo-config \
config set-context foo@kubernetes \
--user=foo \
--cluster=kubernetes \
--namespace=default
Context "foo@kubernetes" created
```

4. Establecer el contexto de forma predeterminada. Por último se tiene que establecer un contexto de forma predeterminada para acceder al cluster.

```
kubectl --kubeconfig=rbac/foo-config \
config use-context foo@kubernetes
Switched to context "foo@kubernetes".
```


El nuevo fichero `KUBECONFIG` ha quedado configurado correctamente. Consulte su información a través del siguiente comando:

```
kubectl --kubeconfig=rbac/foo-config config view
```

Debe haber notado que para utilizar un fichero de configuración específico será necesario utilizar el parámetro --kubeconfig. Si se incluye este nuevo parámetro de forma constante en las líneas de comando podría volverse tedioso el trabajo con Kubectl. Para solucionar este problema se debe crear la variable de entorno KUBECONFIG para especificar la ubicación del fichero que se va a utilizar. Por defecto esta variable no está establecida en la configuración.
Utilice el siguiente comando para utilizar el fichero mmorejon-config de forma predeterminada. 

```
export KUBECONFIG=$PWD/rbac/mmorejon-config
```


Consulte nuevamente la información pero sin el parámetro `--kubeconfig`.

```
kubectl config view
```

Como podrá comprobar Kubectl está utilizando la configuración del usuario mmorejon. Generación automática de certificados y ficheros KUBECONFIG Hasta el momento ha tenido que realizar siete comandos para generar los certificados del usuario y
el fichero KUBECONFIG. Piense que si por alguna casualidad o error se ha equivocado en alguno de estos pasos tendrá que repetir el procedimiento desde el inicio. Para evitar posibles errores y garantizar obtener siempre el mismo resultado se ha creado el fichero bash/create-user.sh. Este script tiene como objetivo crear los certificados y el fichero KUBECONFIG
a partir de un nombre de usuario y el grupo al que pertenece. Antes de utilizarlo es recomendable eliminar los ficheros creados para el usuario mmorejon. Ahora puede crear los ficheros correspondientes al usuario mmorejon perteneciente al grupo developer.

### Comprobar los permisos del nuevo usuario

Una vez que ha quedado configurado el contexto del usuario mmorejon se debe comprobar si este puede
listar los Pods del Namespace default. 

```
kubectl get pods
Error from server (Forbidden): pods is forbidden: User "mmorejon" cannot list resource "p\
ods" in API group "" in the namespace "default"
```

Al ejecutar el comando va a recibir un mensaje de error donde le dice que el usuario mmorejon no tiene permisos para realizar esta operación. Este mensaje es un buen síntoma porque la petición ha pasado correctamente la fase de Autenticación, pero por otro lado es malo porque los permisos del usuario no son suficientes. Para establecer correctamente los permisos del usuario deberá conocer cómo Kubernetes gestiona la fase Autorización.
Antes de pasar a la próxima sección debe eliminar la variable KUBECONFIG a través del siguiente
comando. 

```
unset KUBECONFIG
```
Una vez borrada la variable de entorno Kubectl volverá a utilizar el fichero creado por el cluster.

```
kubectl get pods
No resources found in default namespace.
```

### Roles y RoleBindings

La fase de Autorización será la encargada de validar las peticiones. La validación será realizada en base a los permisos que tiene asignado el usuario o servicio que desea realizar la acción. Lea nuevamente el mensaje de error que recibió el usuario mmorejon: 

```
Error from server (Forbidden): pods is forbidden: User “mmorejon” cannot list resource “pods” in API
group “” in the namespace “default”
```

En este mensaje se pueden apreciar los componentes que son importantes para validar las peticiones: el usuario (mmorejon), la acción a realizar (list), el recurso a ser utilizado (pods), el API a la que pertenece este recurso (“”) y por último el Namespace donde será realizada la acción (default). Estos elementos serán los que formen parte del mecanismo de autorización que ha establecido Kubernetes para acceder a los objetos, y se conoce como Acceso Basado en Roles.
El Control de Acceso Basado en Roles (RBAC) es la configuración predeterminada que tiene el cluster para validar todas las peticiones que llegan al API. Este sistema tiene al Rol como estructura principal que define las acciones que pueden ser realizadas sobre los objetos. En un cluster pueden haber tantos Roles como sean necesarios para cubrir las necesidades de accesos a los recursos, luego los roles son asociados a los usuarios para que estos últimos puedan realizar las operaciones deseadas. A través del RBAC se garantiza una clara separación entre usuarios y responsabilidades. Si en algún momento es necesario modificar o revocar los permisos de un usuario se podrá hacer de forma fácil y sencilla a través de la gestión de roles.

### Roles

Los roles son objetos dentro del cluster de tipo Role y su estructura es la siguiente.