---
id: pod
title: Pods
sidebar_label: Pods
sidebar_position: 3
---

# Pod

:::info

Todos los comandos se ejecutan en el namespace de "lab" para evitar tener que indicar en cada comando el argumento `--namespace lab` podemos cambiar el contexto de ejecución del namespace por defecto (`default`) al namespace de `lab` mediante:

```
kubectl config set-context --current --namespace lab
```

:::

## Introducción

La unidad más pequeña de ejecución que puede utilizar Kubernetes es el [*Pod*](https://kubernetes.io/es/docs/concepts/workloads/pods/pod/), en inglés Pod significa "vaina", y podemos entender un Pod como una envoltura que contiene uno o varios contenedores (en la mayoría de los casos un solo contenedor). Las principales características que tiene un pod son:

![](./02/img/02-pods.png#center)

+ Encapsula a uno o varios contenedores en ejecución. 
+ Tiene asignada una IP, todos los contenedores dentro del Pod comparten la misma IP y puertos. 
+ Los contenedores dentro de un Pod se pueden comunicar entre si utilizando `localhost`. 

![](./02/img/03-pods.png#center)

## Pod Creation Flow

## Maneras de crear objetos en K8s

### Imperativa

```
$ kubectl create namespace hola
$ kubectl run nginx --image=nginx -n hola
$ kubectl edit pod/nginx -n hola
```

### Declarativa

```
$ vim ngix-pod.yaml
$ kubectl apply -f ngix-pod.yaml
$ kubectl delete -f ngix-pod.yaml
```

### Hybrida

```
$ kubectl run nginx --image=nginx --dry-run=client -o yaml > ngix-pod.yaml
$ vim ngix-pod.yaml
$ kubectl create -f ngix-pod.yaml
```





## Crear un pod a partir de un fichero YAML


La estructura básica para crear un Pod podría ser el contenido del fichero [`pod.yaml`](files/pod.yaml):

```yaml
apiVersion: v1 # required
kind: Pod # required
metadata: # required
 name: pod-nginx # required
 labels:
   app: nginx
   service: web
spec: # required
 containers:
   - image: nginx:1.16
     name: contenedor-nginx
     imagePullPolicy: Always
```

 Veamos cada uno de los parámetros que hemos definido:

* `apiVersion: v1`: La versión de la API que vamos a usar.
* `kind: Pod`: La clase de recurso que estamos definiendo.
* `metadata`: Información que nos permite identificar unívocamente el recurso:
    * `name`: Nombre del pod
    * `labels`: Las [Labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) nos permiten etiquetar los recursos de Kubernetes (por ejemplo un pod) con información del tipo clave/valor.
* `spec`: Definimos las características del recurso. En el caso de un Pod indicamos los contenedores que van a formar el Pod (sección `containers`), en este caso sólo uno.
    * `image`: La imagen desde la que se va a crear el contenedor
    * `name`: Nombre del contenedor.
    * `imagePullPolicy`: Las imágenes se guardan en un registro interno. Se pueden utilizar registros públicos (google o docker hub son los más usados) y registros privados. La política por defecto es `IfNotPresent`, que se baja la imagen si no está en el registro interno. Si queremos forzar la descarga desde el repositorio externo, tendremos que indicar `imagePullPolicy:Always`.


Un Pod tiene otros muchos parámetros asociados, que en este caso quedarán sin definir o Kubernetes asumirá los valores por defecto.

Utilizando la forma declarativa podemos crear el pod ([`pod.yaml`](files/pod.yaml)) que hemos visto en la sección anterior :

```
$ kubectl create -f pod.yaml 
pod/pod-nginx created
```

El comando `kubectl create -f` se puede usar para crear cualquier recurso (no solamente pods) a partir de un fichero YAML o JSON. 

:::info
El Pod creado anteriormente de forma declarativa lo podemos crear de manera imperativa mediante `kubectl` de la siguiente manera:

```
$ kubectl run pod-nginx --image=nginx
```

De esta forma se crea un Pod con un contenedor que utiliza la imagen `nginx:latest` (no hemos especificado una versión) del registro que esté definido por defecto en el cluster de Kubernetes, se asigna una dirección IP y se lanza en uno de los nodos del cluster. 
:::

#### Listando los Pods 

El pod se ha creado pero, ¿cómo sabemos que esta corriendo? Vamos a listar los pods que tenemos ejecutandose en el cluster y a visualizar su estado:

```
$ kubectl get pods
NAME        READY   STATUS    RESTARTS   AGE
pod-nginx   1/1     Running   0          51s
```

El resultado del comando anterior muestra el Pod creado y su estado. Todas las columnas son más o menos explicativas, quizás la menos intuitiva pudiera ser READY. Esta columna hace referencia al estado de los contenedores definidos en el Pod, __contenedores listos / total__. En el caso del ejemplo, la columna READY nos indica que se está ejecutando un contenedor del total que es uno. 

La cinco columnas mostradas en el listado corresponden con la información básica del Pod. Si deseamos ampliar la información podemos utilizar el parámetro `--output wide` o `-o wide`

Si queremos ver más información sobre los Pods, como por ejemplo, saber en qué nodo del cluster se está ejecutando: 

```
$ kubectl get pod -o wide
NAME        READY   STATUS    RESTARTS   AGE   IP           NODE                 
pod-nginx   1/1     Running   0          80s   10.244.0.5   kind-control-plane
```

#### Logs dentro del Pod

Los logs son una de las herramientas que tenemos para poder hacer depuración de errores. Podemos obtener los logs del Pod creado a través del comando:

```
$ kubectl logs pod-nginx
```

Si queremos mantener el seguimiento de los logs del Pod podemos usar el parámetro `--follow` o `-f`

```
kubectl logs pod-nginx -f
```

:::info
Los Pods son los únicos objetos en Kubernetes que utilizan el comando `kubectl logs`
:::

## Obtener información del Pod

La información detallada del Pod puede obtenerse a través del siguiente comando (equivalente al inspect de docker):

```
$ kubectl describe pod pod-nginx
Name:             pod-nginx
Namespace:        lab
Priority:         0
Service Account:  default
Node:             kind-control-plane/172.18.0.2
Start Time:       Wed, 12 Oct 2022 08:24:41 +0200
Labels:           app=nginx
                  service=web
Annotations:      <none>
Status:           Running
IP:               10.244.0.5
IPs:
  IP:  10.244.0.5
Containers:
  contenedor-nginx:
    Container ID:   containerd://96cd761227d50e5bd2a1f491d619810650fa7c4ea165d863d5c14d7fac95fc12
    Image:          nginx:1.16
    Image ID:       docker.io/library/nginx@sha256:d20aa6d1cae56fd17cd458f4807e0de462caf2336f0b70b5eeb69fcaaf30dd9c
    Port:           <none>
    Host Port:      <none>
    State:          Running
      Started:      Wed, 12 Oct 2022 08:24:50 +0200
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-dqz24 (ro)
Conditions:
  Type              Status
  Initialized       True
  Ready             True
  ContainersReady   True
  PodScheduled      True
Volumes:
  kube-api-access-dqz24:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:                      <none>
```

Como podemos ver, la información que obtenemos es bastante detallada. Es normal que en este momento no identifiquemos todos los elementos que aparecen, no hay problema poco a poco lograremos descrifrar  

### Eventos

Al final de los datos mostrados por el comando anterior está la sección __Events__. Los eventos son el conjunto de pasos que ha realizado Kubernetes para poner en funcionamiento nuestro despliegue. Si nos encontramos ante una situación en la que hemos de averiguar porqué nuestro despliegue no funciona, obtener los eventos deberia ser una de las primeras cosas que deberiamos de hacer. 

### Estados

El estado de un Pod en Kubernetes viene representado por un objeto denominado __PodStatus__, dicho objeto, tiene un campo que podemos consultar, llamado __phase__. 

La fase de un Pod es simple, se trata de un resumen de alto nivel donde se indica en la fase en la que se encuentra el Pod. 

+ __Pending__

    + Indica que el Pod ha sido aceptado por Kubernetes, pero uno o más de sus contenedores aún no han sido creados.
    + Esta fase incluye el tiempo que Kubernetes tarda en llevar a cabo la programación con el Kube Scheduler, así como el tiempo que necesita para descargar imágenes a través de la red.
    + Esta fase podría llevar cierto tiempo en ser completada.
    
+ __Running__

    + Indica que el Pod ha sido vinculado a un nodo del clúster específico
    + En esta fase, todos los contenedores que pudiera contener un Pod han sido creados
    + Esta fase también indica que al menos un contenedor se está ejecutando o está en proceso de inicio/reinicio

+ __Succeeded__

    + Indica que todos los contenedores del Pod han arrancado de forma correcta y no es necesario ningún reinicio de los mismos

+ __Failed__

    + Indica que todos los contenedores del Pod han finalizado su ejecución
    + Al menos un contenedor del Pod ha finalizado con un fallo
    + En términos Docker - Sistema Operativo, el contendor en cuestión que falla finalizó con un código de error distinto de 0 por parte del sistema operativo

+ __Unknown__

    + Indica que por alguna razón no se puede obtener el estado del Pod
    + Si en alguna ocasión encontramos esta fase en un Pod, es síntoma de que hay un error en la comunicación del propio host con el Pod

Los Pods realmente nunca se paran y se inician, sólo se crean y se destruyen Kubernetes adopta la misma filosofía que el uso de Docker puro, sale más barato
destruir y crear de nuevo :)

Posíblemente observaremos otros estados del Pod, que no forman explícitamente parte del ciclo de vida del Pod, pero indicarían ciertas operaciones que se están llevando a cabo

+ __Terminating__

    + El Pod está siendo eliminado

+ __ContainerCreating__

    + Indica que el contenedor dentro del Pod está siendo creado


```
$ kubectl describe pods nginx | grep Status:
Status: Running
```

```
$ kubectl get pods nginx -o yaml
```

```
$ kubectl get pod pod-nginx --output custom-columns=NAME:metadata.name,STATUS:status.phase,NODE_IP:status.hostIP,POD_IP:status.
podIP
NAME        STATUS    NODE_IP      POD_IP
pod-nginx   Running   172.18.0.2   10.244.0.5
```

## Ejecutar comandos dentro del contenedor 

Puede haber situaciones en las que queramos entrar dentro del contenedor y explorar el sistema de archivos. Al igual que en Docker, Kubectl permite acceder a los contenedores de un Pod a través del comando `exec` 

Por ejemplo, podriamos utilizar el siguiente comando para acceder al Pod creado y comprobar el valor de las variables de entorno

```
kubectl exec -it pod-nginx -- /bin/sh
```

Observar que no hemos tenido que especificar el tipo de recurso, esto es debido a que el comando `exec` solo esta disponible para el objeto de tipo Pod. 

También es posible ejecutar solamente un comando dentro del contenedor. Por ejemplo, podemos ejecutar el comando `env` para obtener el listado de variables de entorno definidas dentro del pod

```
kubectl exec -it pod-nginx -- env
```


Podemos acceder a la aplicación, redirigiendo un puerto de localhost al puerto de la aplicación:

```
kubectl port-forward pod-nginx 8080:80
```

Y accedemos al servidor web en la url http://localhost:8080.

**NOTA**: Esta no es la forma con la que accedemos a las aplicaciones en Kubernetes. Para el acceso a las aplicaciones usaremos un recurso llamado Service. Con la anterior instrucción lo que estamos haciendo es una redirección desde localhost el puerto 8080 al puerto 80 del Pod y es útil para pequeñas pruebas de funcionamiento, nunca para acceso real a un servicio.

**NOTA2**: El `port-forward` no es igual a la redirección de puertos de docker, ya que en este caso la redirección de puertos se hace en el
equipo que ejecuta `kubectl`, no en el equipo que ejecuta los Pods o los contenedores.


## Eliminar el pod creado

Para eliminar el pod que hemos creado podemos utilizar el siguiente comando

```
$ kubectl delete pod -f pod.yaml
```

o bien especificar directamente el nombre del Pod que queremos eliminar

```
kubectl delete pod pod-nginx
```

## Comunicaciones

+ Cada pod puede hablar con cualquier pod de la red usando las direcciones IP; no es necesario realizar ningun tipo de NAT
+ Por defecto, dentro de un pod, las comunicaciones entre contenedores se realizan a través de localhost. 

![](./02/img/041-pods.png#center)


## Configurar Pods

### Definir variables de entorno para el contenedor

Cuando creamos un Pod, podemos setear las variables de entorno para los contenedores que se ejecutan dentro del Pod. Para congfigurar las variables de entorno incluimos los campos `env` o `envFrom` en el fichero de configuración. 

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: envar-demo
  labels:
    purpose: demonstrate-envars
spec:
  containers:
  - name: envar-demo-container
    image: gcr.io/google-samples/node-hello:1.0
    env:
    - name: DEMO_GREETING
      value: "Hello from the environment"
    - name: DEMO_FAREWELL
      value: "Such a sweet sorrow"
```


1. Creamos el pod con el manifiesto anterior

```
kubectl apply -f https://k8s.io/examples/pods/inject/envars.yaml
```

2. Mostramos los pods en ejecución


```
kubectl get pods -l purpose=demonstrate-envars
```

La salida es similar a 

```
NAME            READY     STATUS    RESTARTS   AGE
envar-demo      1/1       Running   0          9s
```

3. Obtenemos las variables de entorno definidas dentro del Pod

```
kubectl exec envar-demo -- printenv
```


La salida es similar a esta

```
NODE_VERSION=4.4.2
EXAMPLE_SERVICE_PORT_8080_TCP_ADDR=10.3.245.237
HOSTNAME=envar-demo
...
DEMO_GREETING=Hello from the environment
DEMO_FAREWELL=Such a sweet sorrow
```

### Definir un comando con argumentos

La mayoria de imagenes de los contenedores tienen definido un `ENTRYPOINT` o una instrucción `CMD`. El comando asignado a la instrucción `CMD` es ejecutado como parte del proceso de inicio del contenedor. Podemos redefinir en la definición del Pod tanto el `ENTRYPOINT` como el `CMD` de la imagen o bien especificar el comando a ejecutar por el contenedor en el caso de que no este especificado en la imagen. 


| Description                           | Docker field name | Kubernetes field name |
|---------------------------------------|-------------------|-----------------------|
| The command run by the container      | Entrypoint        | command               |
| The arguments passed to the command   | Cmd               | args                  |


```yaml
apiVersion: v1
kind: Pod
metadata:
  name: command-demo
  labels:
    purpose: demonstrate-command
spec:
  containers:
  - name: command-demo-container
    image: debian
    command: ["printenv"]
    args: ["HOSTNAME", "KUBERNETES_PORT"]
  restartPolicy: OnFailure
```

1. Creamos un pod a partir del manifiesto anterior

```
kubectl apply -f https://k8s.io/examples/pods/commands.yaml
```

2. Obtenemos los pods en ejecución

```
kubectl get pods
```

3. La salida muestra que el contenedor que se ejecuta en el pod command-demo ha finalizado. Para ver 
la salida del comando que ha ejecutado el contenedor podemos visualizar los logs del Pod

```
kubectl logs command-demo
```

La salida muestra los valores de las variables de entorno HOSTNAME y KUBERNETES_PORT:

```
command-demo
tcp://10.3.240.1:443
```

## Pod con un solo contenedor

En la situación más habitual, se definirá un Pod con un contenedor en su interior para ejecutar un solo proceso y este Pod estará ejecutándose mientras lo haga el correspondiente proceso dentro del
contenedor. Algunos ejemplos pueden ser: ejecución en modo demonio de un servidor web, ejecución de un servidor de aplicaciones Java, ejecución de una tarea programada, ejecución en modo demonio de un servidor DNS, etc.

## Pod multicontenedor

En algunos casos la ejecución de un solo proceso por contenedor no es la solución ideal, ya que existen procesos fuertemente acoplados que no pueden comunicarse entre sí fácilmente si se ejecutan en diferentes sistemas, por lo que la solución planteada en esos casos es definir un
Pod multicontenedor y ejecutar cada proceso en un contenedor, pero que puedan comunicarse entre sí como si lo estuvieran haciendo en el mismo sistema, utilizando un dispositivo de almacenamiento compartido si hiciese falta (para leer, escribir ficheros entre ellos) y compartiendo externamente una misma dirección IP. Un ejemplo típico de un Pod multicontenedor es un servidor web nginx con un servidor de aplicaciones PHP-FPM, que se implementaría mediante un solo Pod, pero
ejecutando un proceso de nginx en un contenedor y otro proceso de php-fpm en otro contenedor.

![](./02/img/pod-4.png#center)

### Implicaciones de múltiples contenedores en un pod

Todos los contenedores dentro de un pod comparten un mismo namespace de red (`netns`) y por lo tanto las interfaces de red, las direccion(es) IP y los  puertos que pertenecen a este. 

Debido a que los puertos son compartidos, los contenedores tienen que usar puertos diferentes para exponer sus servicios. Un mismo puerto no puede ser usado al mismo tiempo por dos o más contenedores. Por ejemplo, si levantaramos 2 contenedores de `nginx` en el mismo pod obtendriamos un error al intentar desplegarlo. 

![](./02/img/pod-5.png#center)

:::info
Cuando desplegamos un pod con más de un container, tenemos que tener la
precaución de que operen en puertos distintos.
:::


A nivel de pod, todos los contenedores comparten los mismos recursos (cgroups, namespaces, volúmenes, etc.)