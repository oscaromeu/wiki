---
id: objetos
title: Objetos
sidebar_label: Objetos
sidebar_position: 2
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

## Gestión imperativa

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