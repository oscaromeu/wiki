---
id: pod
title: Pods
sidebar_label: Pods
sidebar_position: 3
---

# Pod

## Introducción

La unidad más pequeña de ejecución que puede utilizar Kubernetes es el [*Pod*](https://kubernetes.io/es/docs/concepts/workloads/pods/pod/), en inglés Pod significa "vaina", y podemos entender un Pod como una envoltura que contiene uno o varios contenedores (en la mayoría de los casos un solo contenedor). Las principales características que tiene un pod son:

![](./03/img/pod-1.png#center)

+ Encapsula a uno o varios contenedores en ejecución. 
+ Tiene asignada una IP, todos los contenedores dentro del Pod comparten la misma IP y puertos. 
+ Los contenedores dentro de un Pod se pueden comunicar entre si utilizando `localhost`. 

![](./03/img/pod-2.png#center)



## Describiendo un Pod

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

## Gestionando los Pods

### Creando pods usando ficheros YAML

Vamos a crear un pod a partir del fichero ([`pod.yaml`](files/pod.yaml)) que hemos visto en la sección anterior :

```yaml
apiVersion: v1
kind: Pod
metadata:
 name: pod-nginx
 labels:
   app: nginx
   service: web
spec:
 containers:
   - image: nginx:1.16
     name: contenedor-nginx
     imagePullPolicy: Always
```

#### Crear un pod mediante `kubectl create`

Para crear el pod a partir del fichero anterior ejecutamos

```
$ kubectl create -f pod.yaml 
pod/pod-nginx created
```

El comando `kubectl create -f` se puede usar para crear cualquier recurso (no solamente pods) a partir de un fichero YAML o JSON. 

:::info
Es posible crear un Pod de manera imperativa mediante `kubectl` de la siguiente manera:

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

#### Obtener información del Pod

Para obtener información más detallada del Pod (equivalente al inspect de docker):

```
kubectl describe pod pod-nginx
```

#### Eventos

#### Estados


#### Acceder al contenedor creado

Al igual que en Docker, Kubectl permite acceder a los contenedores de un Pod a través del comando `exec` 

Por ejemplo, podriamos utilizar el siguiente comando para acceder al Pod creado y comprobar el valor de las variables de entorno

```
    kubectl exec -it pod-nginx -- /bin/bash
```

Podemos acceder a la aplicación, redirigiendo un puerto de localhost
al puerto de la aplicación:

    kubectl port-forward pod-nginx 8080:80

Y accedemos al servidor web en la url http://localhost:8080.

**NOTA**: Esta no es la forma con la que accedemos a las aplicaciones en Kubernetes. Para el acceso a las aplicaciones usaremos un recurso llamado Service. Con la anterior instrucción lo que estamos haciendo es una redirección desde localhost el puerto 8080 al puerto 80 del Pod y es útil para pequeñas pruebas de funcionamiento, nunca para acceso real a un servicio.
**NOTA2**: El `port-forward` no es igual a la redirección de puertos
de docker, ya que en este caso la redirección de puertos se hace en el
equipo que ejecuta `kubectl`, no en el equipo que ejecuta los Pods o
los contenedores.

Para obtener las etiquetas de los Pods que hemos creado:

    kubectl get pods --show-labels

Las etiquetas las hemos definido en la sección metadata del fichero
yaml, pero también podemos añadirlos a los Pods ya creados:

    kubectl label pods pod-nginx service=web --overwrite=true

Las etiquetas son muy útiles, ya que permiten seleccionar un recurso determinado (en un cluster de Kubernetes puede haber cientos o miles de objetos).Por ejemplo para visualizar los Pods que tienen una etiqueta con un determinado valor:

    kubectl get pods -l service=web

También podemos visualizar los valores de las etiquetas como una nueva
columna:

    kubectl get pods -Lservice

Y por último, eliminamos el Pod mediante:

    kubectl delete pod pod-nginx

Un aspecto muy importante que hay que ir asumiendo es que los Pods son efímeros, se lanzan y en determinadas circunstancias se paran o se destruyen, creando en muchos casos nuevos Pods que sustituyan a los anteriores. Esto tiene importantes ventajas a la hora de realizar modificaciones en los despliegues en producción, pero tiene una consecuencia directa sobre la información que pueda tener almacenada el Pod, por lo que tendremos que utilizar algún mecanismo adicional cuando necesitemos que la información sobreviva a un Pod. Por lo tanto, aunque Kubernetes es un orquestador de contenedores, **la unidad mínima de ejecución es el Pod**, que contendrá uno a más contenedores según las necesidades: 

* En la mayoría de los casos y siguiendo el principio de un proceso por contenedor, evitamos tener sistemas (como máquinas virtuales) ejecutando docenas de procesos, por lo que lo más habitual será tener un Pod en cuyo interior se define un contenedor que ejecuta un solo proceso. 


* En determinadas circunstancias será necesario ejecutar más de un proceso en el mismo "sistema", como en los casos de procesos fuertemente acoplados, en esos casos, tendremos más de un contenedor dentro del Pod. Cada uno de los contenedores ejecutando un solo proceso, pero pudiendo compartir almacenamiento y una misma dirección IP como si se tratase de un sistema ejecutando múltiples procesos.

![](./03/img/pod-3.png#center)

Existen además algunas razones que hacen que sea conveniente tener esta capa adicional por encima de la definición de contenedor:

* Kubernetes puede trabajar con distintos sistemas de gestión de contenedores (docker, containerd, rocket, cri-o, etc) por lo que es muy conveniente añadir una capa de abstracción que permita utilizar Kubernetes de una forma homogénea e independiente del sistema de contenedores interno asociado. 
* Esta capa de abstracción añade información adicional necesaria en Kubernetes como por ejemplo, políticas de reinicio, comprobaciones de que la aplicación esté inicializada (readiness probe) o comprobaciones de que la aplicación haya realizado alguna acción especificada (liveness probe).

## Pod con un solo contenedor

En la situación más habitual, se definirá un Pod con un contenedor en su interior para ejecutar un solo proceso y este Pod estará ejecutándose mientras lo haga el correspondiente proceso dentro del
contenedor. Algunos ejemplos pueden ser: ejecución en modo demonio de un servidor web, ejecución de un servidor de aplicaciones Java, ejecución de una tarea programada, ejecución en modo demonio de un servidor DNS, etc.

## Pod multicontenedor

En algunos casos la ejecución de un solo proceso por contenedor no es la solución ideal, ya que existen procesos fuertemente acoplados que no pueden comunicarse entre sí fácilmente si se ejecutan en diferentes sistemas, por lo que la solución planteada en esos casos es definir un
Pod multicontenedor y ejecutar cada proceso en un contenedor, pero que puedan comunicarse entre sí como si lo estuvieran haciendo en el mismo sistema, utilizando un dispositivo de almacenamiento compartido si hiciese falta (para leer, escribir ficheros entre ellos) y compartiendo externamente una misma dirección IP. Un ejemplo típico de un Pod multicontenedor es un servidor web nginx con un servidor de aplicaciones PHP-FPM, que se implementaría mediante un solo Pod, pero
ejecutando un proceso de nginx en un contenedor y otro proceso de php-fpm en otro contenedor.

![](./03/img/pod-4.png#center)

### Implicaciones de múltiples contenedores en un pod

Todos los contenedores dentro de un pod comparten un mismo namespace de red (`netns`) y por lo tanto las interfaces de red, las direccion(es) IP y los  puertos que pertenecen a este. 

Debido a que los puertos son compartidos, los contenedores tienen que usar puertos diferentes para exponer sus servicios. Un mismo puerto no puede ser usado al mismo tiempo por dos o más contenedores. Por ejemplo, si levantaramos 2 contenedores de `nginx` en el mismo pod obtendriamos un error al intentar desplegarlo. 

![](./03/img/pod-5.png#center)

:::info
Cuando desplegamos un pod con más de un container, tenemos que tener la
precaución de que operen en puertos distintos.
:::


A nivel de pod, todos los contenedores comparten los mismos recursos (cgroups, namespaces, volúmenes, etc.)