---
id: objetos
title: Objetos
sidebar_label: Objetos
sidebar_position: 2
draft: true
---

Los objetos en Kubernetes son entidades que persisten en el sistema, específicamente en el componente `etcd`. Cada objeto representa un estado en Kubernetes y entre todos conforman el estado del cluster. Algunos ejemplos de las representaciones de estos estados puede ser:

+ las aplicaciones que están en ejecución.
+ los recursos disponibles para una aplicación
+ las políticas de redes de una aplicación.

:::info
Entre los objetos más utilizados están los `Pods`, `Deployments`, `ConfigMaps` y `Secrets`
:::

Cada objeto creado en Kubernetes es un nuevo estado que debe ser establecido por el cluster. Desde el mismo momento de su creación, Kubernetes va a trabajar constantemente para garantizar que este objeto logre alcanzar el estado deseado. Para crear, modificar o eliminar objetos tiene que utilizarse la API de Kubernetes. En este punto es donde se utiliza Kubectl⁴⁹, ya que brinda una interfaz de comandos intuitiva y fácil para interactuar con el cluster.

## Los objetos Spec y Status

Cada objeto en Kubernetes está compuesto por dos campos anidados, donde cada uno es otro objeto. Estos dos objetos anidados son Spec y Status. El objeto Spec describe las características del estado deseado y es proporcionado desde el cliente, p.ej utilizando Kubectl. El objeto Status representa el estado actual dentro del cluster y es actualizado por el sistema Kubernetes. El plano de control es el responsable de observar constantemente los objetos y de lograr que todos alcancen el estado deseado. 

Por ejemplo, si en las características del objeto se establece que la aplicación debe tener tres réplicas, el sistema de Kubernetes va a escalar la aplicación a tres para lograr el estado deseado.

![](./00/img/01-spec-status.png#center)

## Comandos desde Kubectl


Todos los elementos que existen dentro de Kubernetes son representados a través de objetos y gestionados desde un cliente. Entre los clientes más utilizados está Kubectl por brindar una poderosa y simple estructura de comandos para realizar las operaciones en el cluster.
Listado de comandos de kubectl:

```
annotate        completion      drain       options         set
api-resources   config          edit        patch           taint
api-versions    convert         exec        plugin          top
apply           cordon          explain     port-forward    uncordon
attach          cp              expose      proxy           version
auth            create          get         replace         wait
autoscale       delete          kustomize   rollout
certificate     describe        label       run
cluster-info    diff            logs        scale
```

### Familia de comandos

La mayoría de los comandos presentes en Kubectl pueden agruparse en una de las siguientes categorías:

#### Gestión declarativa`

La forma recomendada de gestionar los recursos es a través de los ficheros declarativos llamados Configuración de Recursos (Resource Config) y del comando kubectl apply. Esta gestión declarativa de los recursos se utiliza principalmente para despliegues y operaciones.
En los ficheros Resource Config es donde se especifica el estado deseado del objeto (objeto Spec). La estructura de los ficheros es YAML, p.ej:

```yaml
apiVersion: v1
kind: Pod
metadata:
  namespace: monitoring
  labels:
    run: nginx
    name: nginx
spec:
  containers:
    - image: nginx
      name: nginx
      resources: {}
    dnsPolicy: ClusterFirst
    restartPolicy: Never
```

Luego se utiliza el comando Apply de Kubectl para enviar la información a Kubernetes y lograr el estado deseado (objeto Status).

```
kubectl apply -f nombre-del-fichero.yaml
```

Para crear un objeto a través de un fichero YAML existen cuatro elementos que no pueden faltar:

|Atributo| Descripción|
|--------|------------|
|`apiVersion`| Especifica la versión del API de Kubernetes utilizada para crear el objeto.|
|`kind`| Tipo de objeto que se desea crear.|
|`metadata`| Información para identificar de forma única al objeto.|
|`spec`| Descripción del estado que se desea que obtenga el objeto dentro del cluster.|

### Gestión imperativa

A través de la forma imperativa Kubectl indica la acción que debe realizarse en cada comando. p.ej:
+ `kubectl create deployment <deployment-name>` para crear un objeto Deployment.
+ `kubectl delete pod <pod-name>` para borrar un objeto Pod.
+ `kubectl set resources deployment <deployment-name>` para asignarle recursos a un objeto Deployment.

En cada línea de comando se puede identificar claramente la acción a realizar por los verbos utilizados. Esta es la principal diferencia con la forma declarativa, donde la acción es identificada por Kubectl en dependencia del contenido del fichero YAML.

## Impresión de estados

Es necesario poder ver el estado de los objetos y por tal motivo Kubectl cuenta con comandos para:

+ Mostrar el estado y la descripción de un objeto.
+ Mostrar campos específicos de un objeto.
+ Realizar consultas a los objetos basados en etiquetas

## Interacción con contenedores
Kubectl permite depurar los procesos en ejecución dentro del cluster y para ello cuenta con comandos
para:
+ Imprimir los logs de los contenedores
+ Imprimir los eventos realizados en el cluster
+ Acceder al contenedor para efectuar comandos
+ Copiar ficheros desde la máquina local hacia el contenedor

## Gestión del cluster
Es muy probable que necesite en algún momento realizar operaciones de mantenimiento al cluster. Con este propósito Kubectl brinda comandos para aislar o drenar los nodos, p.ej: drain, uncordon, cordon, entre otros.

## Tipos de objetos disponibles
Son muchos los tipos de objetos que pueden funcionar en su cluster y aprenderse cada uno de los nombres puede tomar tiempo. Siempre que necesite tener información asociada a los tipos de objetos, sus alias y grupo al que pertenece, utilice el comando api-resources. Consulte el siguiente fragmento de información brindada por este comando:

```
$ kubectl api-resources
NAME                              SHORTNAMES   APIVERSION                             NAMESPACED   KIND
bindings                                       v1                                     true         Binding
componentstatuses                 cs           v1                                     false        ComponentStatus
configmaps                        cm           v1                                     true         ConfigMap
endpoints                         ep           v1                                     true         Endpoints
events                            ev           v1                                     true         Event
limitranges                       limits       v1                                     true         LimitRange
namespaces                        ns           v1                                     false        Namespace
nodes                             no           v1                                     false        Node
persistentvolumeclaims            pvc          v1                                     true         PersistentVolumeClaim
persistentvolumes                 pv           v1                                     false        PersistentVolume
pods                              po           v1                                     true         Pod
podtemplates                                   v1                                     true         PodTemplate
replicationcontrollers            rc           v1                                     true         ReplicationController
resourcequotas                    quota        v1                                     true         ResourceQuota
secrets                                        v1                                     true         Secret
serviceaccounts                   sa           v1                                     true         ServiceAccount
services                          svc          v1                                     true         Service
mutatingwebhookconfigurations                  admissionregistration.k8s.io/v1        false        MutatingWebhookConfiguration
validatingwebhookconfigurations                admissionregistration.k8s.io/v1        false        ValidatingWebhookConfiguration
customresourcedefinitions         crd,crds     apiextensions.k8s.io/v1                false        CustomResourceDefinition
apiservices                                    apiregistration.k8s.io/v1              false        APIService
controllerrevisions                            apps/v1                                true         ControllerRevision
...
```

## API fundamentals

### Uses of Kubernetes API

The Kubernetes API is the core of Kubernetes. The API is the remote interface that enables users to create, configure, and delete Kubernetes objects. It allows developers to extend Kubernetes by creating controllers that drive the system. The controllers are responsible for things like replication, scaling, and rollout of new versions.

The Kubernetes API can be used to deploy applications, pull images from the registry, or even monitor cluster state. In addition, it can also be used for performing rolling upgrades and managing network policies.

The API provides programmatic access to all the features in Kubernetes. With it, developers can automate the deployment of applications, manage workloads, or access information about the cluster. In addition, familiarity with the API helps one know how to best use Kubernetes and kubectl commands to manage pods, services, replication controllers, etc., without needlessly typing commands.

Verify the currently available Kubernetes API versions on the cluster:

:::info

```
$ kubectl config view --minify -oyaml | yq '.clusters[].cluster.server'
https://127.0.0.1:34591


curl -k $(kubectl config view --minify -oyaml | yq ".clusters[].cluster.server")/version 

{
  "major": "1",
  "minor": "24",
  "gitVersion": "v1.24.0",
  "gitCommit": "4ce5a8954017644c5420bae81d72b09b735c21f0",
  "gitTreeState": "clean",
  "buildDate": "2022-05-19T15:39:43Z",
  "goVersion": "go1.18.1",
  "compiler": "gc",
  "platform": "linux/amd64"
}
```
:::

Observad que si intentamos obtener un recurso que no sea version llamando a la api mediante curl este nos va a devolver el siguiente mensaje de error

```
$ curl -k $(kubectl config view --minify -oyaml | yq ".clusters[].cluster.server")/api/v1/pods     
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "pods is forbidden: User \"system:anonymous\" cannot list resource \"pods\" in API group \"\" at the cluster scope",
  "reason": "Forbidden",
  "details": {
    "kind": "pods"
  },
  "code": 403
}%
```

Para poder usar curl debemos especificar en la petición los certificados. Por ahora vamos a usar un proxy para poder acceder a la API

```bash
$ kubectl proxy                                                                          
Starting to serve on 127.0.0.1:8001
```

Y con esto ya podemos lanzar peticiones a la API

```
$ curl -s http://localhost:8001/apis -k | grep "name"
      "name": "apiregistration.k8s.io",
      "name": "apps",
      "name": "events.k8s.io",
      "name": "authentication.k8s.io",
      "name": "authorization.k8s.io",
      "name": "autoscaling",
```


```bash
kubectl api-version
```

Use the `--v` flag to set a verbosity level. This will allow you to see the request/responses against the Kubernetes API:

```
kubectl get pods --v=8
```

Use the `kubectl proxy` command to proxy local requests on port 8001 to the Kubernetes API:

```
kubectl proxy --port=8001
```

Open up another terminal by clicking the + button and select Open New Terminal.

Send a GET request to the Kubernetes API using curl:

```
curl -X GET http://localhost:8001
```

You can explore the OpenAPI definition file to see complete API details.

```
curl localhost:8001/openapi/v2
```

Send a GET request to list all pods in the environment:

```
curl -X GET http://localhost:8001/api/v1/pods
```

Use jq to parse the json response:

```
curl -X GET http://localhost:8001/api/v1/pods | jq .items[].metadata.name
```

You can scope the response to a particular namespace:

```
curl -X GET http://localhost:8001/api/v1/namespaces/myproject/pods
```

Get more details on a particular pod within the myproject namespace:

```
curl -X GET http://localhost:8001/api/v1/namespaces/myproject/pods/my-two-container-pod
```

Export the manifest associated with my-two-container-pod in json format:

```
kubectl get pods my-two-container-pod -o json | jq 'del(.metadata.namespace,.metadata.resourceVersion,.metadata.uid) | .metadata.creationTimestamp=null' > podmanifest.json
```

Within the manifest, replace the 1.13 version of alpine with 1.14:

```
sed -i 's|nginx:1.13-alpine|nginx:1.14-alpine|g' podmanifest.json
```

Update/Replace the current pod manifest with the newly updated manifest:

```
curl -X PUT http://localhost:8001/api/v1/namespaces/myproject/pods/my-two-container-pod -H "Content-type: application/json" -d @podmanifest.json
```


You've upgraded the version of alpine in the my-two-container-pod to 1.14. Now, patch the current pod with an even newer container image, 1.15:

```
curl -X PATCH http://localhost:8001/api/v1/namespaces/myproject/pods/my-two-container-pod -H "Content-type: application/strategic-merge-patch+json" -d '{"spec":{"containers":[{"name": "server","image":"nginx:1.15-alpine"}]}}'
```

Delete the current pod by sending the DELETE request method:

```
curl -X DELETE http://localhost:8001/api/v1/namespaces/myproject/pods/my-two-container-pod
```

Verify the pod is in Terminating status by running kubectl get pods. Once it's terminated, you can verify the pod no longer exists:

```
curl -X GET http://localhost:8001/api/v1/namespaces/myproject/pods/my-two-container-pod
```

