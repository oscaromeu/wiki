# Gestionando los Pods

## Creando pods usando ficheros YAML

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

### Crear un pod mediante `kubectl create`

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

### Logs dentro del Pod

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

Para obtener información más detallada del Pod (equivalente al inspect de docker):

```
kubectl describe pod pod-nginx
```

### Eventos

### Estados


## Acceder al contenedor creado

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

