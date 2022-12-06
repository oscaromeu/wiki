---
id: state_persistence
title: Persistencia
sidebar_label: Persistencia
sidebar_position: 10
draft: true
---

## Volúmenes en Kubernetes

En la actulidad la persistencia de datos en un sistema es un tema muy sensible. Toda la información que fluye en un sistema persiste de alguna manera en la infraestructura, ya sea en bases de datos o como ficheros en un un directorio. 

Docker ha asociado el concepto de __Volumen__ para gestionar la información que necessita persistir fuera del ciclo de vida de un contenedor. Usando este mecanismo los contenedores podrán inciarse y detenerse todas las veces que sean necesarias sin que se elimine la información que debe persistir. 

Un volumen corresponde con uno o múltiples ficheros dentro del sistema operativo. Esta estructura de ficheros se incluye dentro del contenedor al ser iniciado y luego permanacen en el sistema operativo al deternerse. 

Kubernetes también incluye el concepto de __Volumen__ para gestionar la persistencia. En Kubernetes un volumen sigue siendo una estructura de ficheros, pero con la diferencia que se incluye primero como parte del Pod y luego se asocia a uno o varios de los contenedores definidos dentro del Pod. 

Esta estructura de ficheros tiene que existir fisicamente en algún componente dentro de la infraestructura. Kubernetes tiene definido un tipo de volumen para casi todas las opciones existentes en las plataformas, p.ej, _emptyDir, hostPath, awsElasticBlockStore, azureDisk, persistentVolumeClaim_ etc.

## Tipos de volúmenes 

Los [volúmenes](https://kubernetes.io/docs/concepts/storage/volumes/) nos permiten proporcionar almacenamiento a los Pods, y podemos usar distintos tipos que nos ofrecen distintas características:

* Proporcionados por proveedores de cloud: AWS, Azure, GCE, OpenStack, etc
* Propios de Kubernetes:
    * configMap, secret: Mecanismo para inyectar datos de configuración a las aplicaciones.
    * emptyDir: Volumen efímero con la misma vida que el Pod. Usado como almacenamiento secundario o para compartir entre contenedores del mismo Pod.
    * hostPath: Monta un directorio del host en el Pod (usado excepcionalmente, pero es el que nosotros vamos a usar con minikube).
    * persistentVolumeClaim: solicitud de uso de almacenamiento.

* Habituales en despliegues "on premises": glusterfs, cephfs, iscsi, nfs, etc.

## Especificaciones del volumen dentro del Pod

Definir un volumen para un pod requiere dos pasos. El primer paso consiste en definir el volumen utilizando el atributo `spec.volumes` Como parte de la definición se debe incluir el nombre y el tipo de volumen. Para usar el volumen no es suficiente con haberlo declarado. Como segundo paso el volumen tiene que ser montado en una ruta para que el contenedor pueda usarlo utilizando el atributo `spec.containers.volumeMounts`. Es obligatorio que el campo `name` corresponda con uno de los elementos listados en la estructura `volumes`. 

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: business-app
spec:
  volumes:
  - name: logs-volume
    emptyDir: {}
  containers:
  - image: nginx
    name: nginx
    volumeMounts:
    - mountPath: /var/logs
      name: logs-volume
```

Notar la definición del volumen `logs-volume` y su correspondiente tipo `emptyDir`. Observad como este se ha montado en la ruta `/var/logs` del contenedor `nginx`.

## Provisionado dinámico o estático de volumenes

Los objetos de tipo `PersistentVolume` pueden ser creados de manera estática o dinámica. Si optamos por el enfoque estático, primero debemos crear un dispositivo de almacenamiento y hacer referencia a él creando explícitamente un objeto de tipo PersitentVolume. El enfoque dinámico no requiere crear un objeto PersitentVolume. Se creará automáticamente a partir del PersitentVolumeClaim estableciendo un nombre de clase de almacenamiento mediante el atributo `spec.storageClassName`. Una clase de almacenamiento es un concepto de abstracción que define una clase de dispositivo de almacenamiento (por ejemplo, almacenamiento con performance lento o rápido). Generar este tipo de configuraciones y/o objetos es tarea del equipo de administración. Podemos inspeccionar las clases de almacenamiento disponibles en el cluster con el comando

```
$ kubectl get storageclass
```

## PersistentVolume

El _Persitent Volume (PV)_ es un objeto en el cluster que representa un almacenamiento, es un volumen persistente y constante. Este elemento es proporcionado por el equipo de administración o de manera dinámica. Los PV son objetos que no están asociados a un Namespace. 

Los principales elementos que forman parte de un _Persisten Volume_ son:

|Campo|Descripción|
|-----|-----------|
|`kind`|El tipo de objeto utilizado, en este caso es PersitentVolume|
|`metadata`|Establece la descripción del objeto a través del nombre y las etiquetas|
|`storage`|Define el máximo de almacenamiento a utilizar.|
|`accessModes`|Establece el tipo de accesso al volumen. Las opciones posibles son: _ReadWriteOnce_, _ReadOnlyMany_ y _ReadWriteMany_.
|`persistentVolumeReclaimPolicy`|Define la política de retención de la información en el volumen. Las opciones disponibles son: _Retain_, _Recycle_, y _Delete_.|

De especial interes es el parámetro `accessModes`:

|Tipo|Descripción|
|----|-----------|
|ReadWriteOnce| Acceso lectura/escritura por un solo nodo|
|ReadOnlyMany|Acceso solo lectura por diferentes nodos|
|ReadWriteMany|Accesso lectura/escritura por varios nodos| 

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: db-pv
spec:
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /data/db
```

```
$ kubectl get pv
NAME    CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM   STORAGECLASS   REASON   AGE
db-pv   1Gi        RWO            Retain           Available                                   2m51s
```

En el listado se muestran los detalles del PV creado. La columna STATUS nos indica el estado actual del volumen. Los diferentes estado posibles coresponden con las siguientes fases:

|Fase|Descripción|
|----|-----------|
|Available| Recurso disponible. No se encuentra atado a ninguna solicitud o demanda|
|Bound|Recurso atado a una solicitud o demanda|
|Released|La solicitud que ataba al recurso ha sido eliminada, pero el recurso no ha sido solicitado por el cluster|
|Failed|El volumen ha presentado un error en su recuperación automática|

## Solicitud de uso del PersistentVolume aka PersitentVolumeClaims

El siguiente objeto que vamos a crear va a ser el `PersitentVolumeClaim` el proposito del cuál es bindear el PersitentVolume con el Pod. El proposito de este paso es bindear el PersitentVolume al Pod. 

```yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: db-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 512m
```

El siguiente comando `get` utiliza la forma corta `pvc` en lugar de `persitentvolumeclaims`

```
$ kubectl create -f db-pvc.yaml
persitentvolumeclaim/db-pvc created
$ kubectl get pvc db-pvc
```

Notar que el PersistentVolume aún no ha sido montado por el Pod, por lo que si inspeccionamos el objeto veremos el estado "None". Utilitzar el comando `describe` es una buena manera de verificar si el PVC ha sido montado.  

```
Access Modes:
VolumeMode:    Filesystem
Used By:       <none>
Events:
  Type    Reason                Age               From                         Message
  ----    ------                ----              ----                         -------
  Normal  WaitForFirstConsumer  8s (x5 over 67s)  persistentvolume-controller  waiting for first consumer to be created before binding
```

## Montar el PersitentVolumeClaims en los pods

El último paso que nos queda por hacer es montar el PersitentVolumeClaim en el Pod que quiere hacer uso de el. Ya hemos visto como montar volumenes en un pod. La gran diferencia aqui es el atributo `spec.persitentVolumeClaim` y que estamos indicando el nombre del PVC que queremos. 

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-consuming-pvc
spec:
  volumes:
    - name: app-storage
      persistentVolumeClaim:
        claimName: db-pvc
  containers:
  - image: alpine
    name: app
    command: ["/bin/sh"]
    args: ["-c", "while true; do sleep 60; done;"]
    volumeMounts:
      - mountPath: "/mnt/data"
        name: app-storage
```

```
$ kubectl get pvc db-pvc
NAME     STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
db-pvc   Bound    pvc-d23c279a-b0c1-49de-872a-692297645943   512m       RWO            standard       7m29s
```

```
$ kubectl describe pvc db-pvc
...
Access Modes:  RWO
VolumeMode:    Filesystem
Used By:       app-consuming-pvc
Events:
  Type    Reason                 Age                     From                                                                                          
      Message
  ----    ------                 ----                    ----                                                                                          
      -------
  Normal  WaitForFirstConsumer   2m11s (x26 over 8m25s)  persistentvolume-controller                                                                   
      waiting for first consumer to be created before binding
  Normal  Provisioning           62s                     rancher.io/local-path_local-path-provisioner-684f458cdd-t6w82_a8a112c3-a939-40ef-a591-3dbfe6b29a60  External provisioner is provisioning volume for claim "default/db-pvc"
  Normal  ProvisioningSucceeded  58s                     rancher.io/local-path_local-path-provisioner-684f458cdd-t6w82_a8a112c3-a939-40ef-a591-3dbfe6b29a60  Successfully provisioned volume pvc-d23c279a-b0c1-49de-872a-692297645943
```