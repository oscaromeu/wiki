---
id: labels
title: Labels y anotations
sidebar_label: Labels y anotations
sidebar_position: 5
draft: true
---

## Introducción

Las etiquetas son pares de clave/valor que se asocian a los objetos, como los pods. El propósito de las etiquetas es permitir identificar atributos de los objetos que son relevantes y significativos para los usuarios, pero que no tienen significado para el sistema principal. Se puede usar las etiquetas para organizar y seleccionar subconjuntos de objetos. Las etiquetas se pueden asociar a los objetos a la hora de crearlos y posteriormente modificarlas o añadir nuevas. Cada objeto puede tener un conjunto de etiquetas clave/valor definidas, donde cada clave debe ser única para un mismo objeto.

![](./04/img/01-labels.png#center)

Ejemplos de etiquetas

+ "release" : "stable", "release" : "canary"
+ "environment" : "dev", "environment" : "qa", "environment" : "production"
+ "tier" : "frontend", "tier" : "backend", "tier" : "cache"
+ "partition" : "customerA", "partition" : "customerB"
+ "track" : "daily", "track" : "weekly"

## Asignar labels

Las etiquetas se pueden definir de manera declarativa en la sección `metadata` de los objetos. Un objeto no tiene restricciones en cuánto al número de etiquetas pero hay que tener en cuenta que no pueden estar duplicadas. 

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: label-demo
  labels:
    env: prod
    app: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
```

También se pueden crear de manera imperativa con el comando `run`

```
$ kubectl run label-demo --image=nginx --restart=Never --labels=env=prod,app=nginx
```

finalmente, también podemos añadir y/o modificar la etiqueta a un pod que esta corriendo dentro del cluster

```
$ kubectl label pod label-demo version=latest
pod/label-demo labeled
```

obtenemos el listado de pods que se estan ejecutando con sus respectivas etiquetas

```
$ kubectl get pods --show-labels
NAME           READY   STATUS    RESTARTS   AGE   LABELS
label-demo              1/1     Running   0          65s   app=nginx,env=prod,version=latest
label-demo-1   1/1     Running   0          25s   app=nginx,env=prod
label-demo-2   1/1     Running   0          23s   app=nginx,env=stagging
label-demo-3   1/1     Running   0          21s   app=nginx,env=prod,version=1.14.2
```

## Filtrar objetos utilizando etiquetas

Podemos seleccionar o filtrar los objetos de kubernets a través de la CLI o bien a través de `spec.selector` en el manifiesto YAML segun sea el caso. La selección se puede realizar con un criterio de igualdad o bien en un criterio basado en pertenencia a un conjunto. Otra manera de obtener las etiquetas asignadas a los objetos de kubernets seria mediante el comando `describe` o `get`. 

__Ejemplo:__ _Obtener las etiquetas de un pod mediante los comandos describe o get_

```
$ kubectl describe pod label-demo | grep -iC 2 labels:
Node:             kind-control-plane/172.18.0.2
Start Time:       Sat, 15 Oct 2022 09:06:49 +0200
Labels:           app=nginx
                  env=prod
                  version=latest
...
$ kubectl get pod label-demo -oyaml | grep -iC 2 labels:
metadata:
  creationTimestamp: "2022-10-15T07:06:49Z"
  labels:
    app: nginx
    env: prod
```


### Requisitos basados en igualdad

![](./04/img/02-labels.png#center)

__Ejemplo:__ _Obtener todos los pods que tienen la etiqueta `env=stagging`_

```
$ kubectl get pods -l env=stagging --show-labels
NAME           READY   STATUS    RESTARTS   AGE   LABELS
label-demo-2   1/1     Running   0          13m   app=nginx,env=stagging
```

Se permiten tres clases de operadores =,==,!=. Los dos primeros representan la igualdad (y son simplemente sinónimos), mientras que el último representa la desigualdad.

__Ejemplo:__ _Obtener todos los pods que tienen la etiqueta `env=prod` pero no su version no es la 1.14.2 (`version!=1.14.2`)_

```
$ kubectl get pods -l env=prod,version!=1.14.2 --show-labels
NAME           READY   STATUS    RESTARTS   AGE   LABELS
label-demo-1   1/1     Running   0          13m   app=nginx,env=prod
```

### Requisitos basados en conjunto

Los requisitos de etiqueta basados en Conjuntos permiten el filtro de claves en base a un conjunto de valores. Se puede utilizar tres tipos de operadores: `in`, `notin` y `exists` (sólo el identificador clave). Por ejemplo:

```
environment in (production, qa)
tier notin (frontend, backend)
partition
!partition
```

El primer ejemplo selecciona todos los recursos cuya clave es igual a environment y su valor es igual a production o qa. El segundo ejemplo selecciona todos los recursos cuya clave es igual a tier y sus valores son distintos de frontend y backend, y todos los recursos que no tengan etiquetas con la clavetier. El tercer ejemplo selecciona todos los recursos que incluyan una etiqueta con la clave partition; sin comprobar los valores. El cuarto ejemplo selecciona todos los recursos que no incluyan una etiqueta con la clave partition; sin comprobar los valores. De forma similar, el separador de coma actúa como un operador AND . Así, el filtro de recursos con una clave igual a partition (sin importar el valor) y con un environment distinto de qa puede expresarse como partition,environment notin (qa). El selector basado en conjunto es una forma genérica de igualdad puesto que environment=production es equivalente a environment in (production); y lo mismo aplica para != y notin.

Los requisitos basados en conjunto pueden alternarse con aquellos basados en igualdad. Por ejemplo: partition in (customerA, customerB),environment!=qa.

__Ejemplo:__ _Obtener todos los pods que tienen definida la etiqueta version sin comprobar su valor_

```
$ kubectl get pods -l version --show-labels
NAME           READY   STATUS    RESTARTS   AGE   LABELS
label-demo-3   1/1     Running   0          24m   app=nginx,env=prod,version=1.14.2
```

## Establecer referencias en los objetos de la API

Algunos objetos de Kubernetes, como los services y los replicationcontrollers, también hacen uso de los selectores de etiqueta para referirse a otros conjuntos de objetos, como los pods.

### Service y ReplicationController
El conjunto de pods que un service expone se define con un selector de etiqueta. De forma similar, el conjunto de pods que un replicationcontroller debería gestionar se define con un selector de etiqueta.

Los selectores de etiqueta para ambos objetos se definen en los archivos json o yaml usando mapas, y únicamente se permite los selectores basados en igualad:

```yaml
selector:
    component: redis
```

este selector (respectivamente en formato json o yaml) es equivalente a `component=redis` o `component in (redis)`.

### Recursos que soportan requisitos basados en conjunto

Algunos recursos más recientes, como el Job, el Deployment, el Replica Set, y el Daemon Set, sí permiten requisitos basados en conjunto.

```yaml
selector:
  matchLabels:
    component: redis
  matchExpressions:
    - {key: tier, operator: In, values: [cache]}
    - {key: environment, operator: NotIn, values: [dev]}
```

matchLabels es un mapa de pares `{key,value}`. Una única combinación de `{key,value}` en el mapa `matchLabels` es equivalente a un elemento en `matchExpressions` donde el campo key es "key", el operator es "In", y la matriz values contiene únicamente "value". matchExpressions es una lista de requisitos de selección de pod. Los operadores permitidos son In, NotIn, Exists, y DoesNotExist. El conjunto de valores no puede ser vacío en el caso particular de In y NotIn. Todos los requisitos, tanto de matchLabels como de matchExpressions se combinan entre sí con el operador AND -- todos ellos deben ser satisfechos.

Seleccionar conjuntos de objetos
Un caso de uso de selección basada en etiquetas es la posibilidad de limitar los nodos en los que un pod puede desplegarse.

## Recomendaciones sobre etiquetas

Es interesante siempre incluir etiquetas en los objetos de Kubernetes. A priori puede parecer que es poco interesante pero con el paso del tiempo el número de aplicaciones desplegadas en el cluster aumentará. Las etiquetas también nos resultaran de utilidad cuando incluyamos en nuestro cluster elementos de monitorización. Las etiquetas nos seran utiles para generar dashboards avanzados segun las diferentes etiquetas que tienen nuestros recursos. 

Otro aspecto a considerar es que se pueden dar situaciones en las que nos pueda interesar restringir el despliegue de aplicaciones en función de alguna característica del cluster. Podemos utilizar el siguiente comando para ver cómo kubernetes establece sus etiquetas

```
kubectl get nodes --show-labels
````

## Anotaciones

Las anotaciones de Kubernetes permiten adjuntar metadatos arbitrarios a los objetos, de tal forma que clientes como herramientas y librerías puedan obtener fácilmente dichos metadatos. Las anotaciones no se utilizan para identificar y seleccionar objetos. Los metadatos de una anotación pueden ser pequeños o grandes, estructurados o no estructurados, y pueden incluir caracteres no permitidos en las etiquetas.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: label-demo
  annotations:
    branch: master
  labels:
    env: prod
    app: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
```