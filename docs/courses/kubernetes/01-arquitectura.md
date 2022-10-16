---
id: arquitectura
title: Arquitectura de Kubernetes
sidebar_label: Arquitectura de Kubernetes
sidebar_position: 1
---

# Arquitectura de Kubernetes

Es importante conocer las piezas que hacen posible el funcionamiento del cluster de Kubernetes para poder detectar, entender y corregir los problemas que surjan en las actividades diarias. 

Cuando se despliega Kubernetes obtendremos como resultado un cluster. Un cluster es un grupo de máquinas llamadas _nodos_ con aplicaciones encapsuladas en contenedores en funcionamiento. Un cluster está compuesto por al menos un nodo, pero es recomendable separar las responsabilidades en diferentes nodos. 

* Los nodos *master*: Son los que ejecutan los servicios principales de k8s y coordinan las actividades del cluster. El uso del término master  es últimamente muy controvertido en los paises de habla inglesa, se está cambiando su denominación por *control plane* node. 

* Los nodos *worker*: Son los que reciben las órdenes de los controladores y en los que se ejecutan los contenedores de las
aplicaciones. 


k8s es un software que se instala en varios nodos que se gestionan de forma coordinada, es decir, un clúster de nodos. Aunque es posible en
casos muy peculiares instalar algunos [nodos sobre sistemas Windows](https://kubernetes.io/docs/setup/production-environment/windows/intro-windows-in-kubernetes/), la situación normal es que se trate de un cluster de nodos linux. No es necesario que todos los nodos tengan la misma versión y ni siquiera que sean la misma distribución, aunque en muchos casos sí lo sea por simplicidad en el despliegue y mantenimiento.

Los nodos del clúster pueden ser máquinas físicas o virtuales, pero quizás lo más habitual es que se traten de instancias de nube de
infraestructura, es decir, máquinas virtuales ejecutándose en algún proveedor de IaaS (AWS, GCP, OpenStack, etc.)

![](./01/img/01-architecture.png#center)

## Componentes de un nodo master



### kube-apiserver

El componente [kube-apiserver] expone el API de kubernetes; es la interfaz de comunicación para acceder al plano de control. Todas las peticiones realizadas a Kubernetes pasan por esta interfaz donde los objetos son examinados y validados antes de ser aplicados.

### etcd

Almacenamiento de llave - valor consistente y de alta disponibilidad utilizado en Kubernetes para almacenar toda la información del cluster

### kube-scheduler

El componente kube-scheduler se mantiene observando el apiserver en busca de nuevos trabajos. Si detecta que ha entrado un nuevo trabajo su misión consiste en seleccionar uno de los nodos existentes y asignar el nuevo trabajo

### kube-controller-manager

El componente kube-controller-manager es un binario compuesto por diferentes controladores de Kubernetes, entre los que se encuentran:

+ Node Controller: Responsable de notificar y responder cuando un nodo deja de funcionar.
+ Replication Controller: Responsable de mantener el número correcto de Pods para cada una de las réplicas.
+ Endpoints Controller: Responsable de poblar los obetos Endpoints, lo que significa enlazar los Pods con los servicios.
+ Service Account & Token Controllers: Responsable de crear las cuentas predeterminadas y los tokens de acceso en los nuevos namespaces.


### cloud-controller-manager

Cloud controller manager

Este controlador se ejecuta cuando estamos ejecutando el cluster en un cloud ya sea AWS, Azure, GCP etc. Su función es gestionar la integración con los servicios o tecnologías dadas por el proveedor. Por ejemplo, si desplegamos una aplicación que solicita por un balanceador expuesto a internet, el cloud controller manager sera el responsable de provisionar dicho balanceador. 


## Componentes de un nodo worker (Data Plane)

Los worker nodes o simplemente nodes a alto nivel ejecutan 3 tareas

+ Observar nuevas cargas de trabajo a través de la API server
+ Ejecutar las cargas de trabajo asignadas
+ Reportar o dar feedback al plano de control via la API server

Cada nodo del data plane está formado en general por estos 3 componentes (puede haber más)


### kubelet 

kubelet es un agente ejecutándose en cada nodo del cluster. Su misión es garantizar que cada contenedor esté funcionando dentro de un Pod. Kubelet no va a gestionar los contenedores que no sean gestionados por Kubernetes

### kube-proxy

Kube proxy mantiene las reglas de red en los nodos, estas reglas permiten las comunicaciones hacia los pods tanto desde afuera como dentro del cluster. Kube proxy implementa parte del objeto service en kubernetes. Kube proxy utilitza utilidades de filtrado de paquetes a nivel de sistema operativo o bien hace el reenvío el mismo. Más sobre esto después. 

### Container Runtime

Kubelet necesita un motor de contenedores. Es decir, un controlador que gestione los contenedores, ya sea la descarga de imagenes, el inicio o parada de los mismos etc. Como todo en kubernetes es modular, es decir, si no nos gusta uno ponemos otro y si eso tampoco nos vale nos hacemos uno nosotros. Dicho de manera más fina, kubernetes soporta cualquier motor de contenedores que cumpla el protocolo principal de comunicación entre kubelet y el motor de contenedores. Este protocolo de comunicaciones esta definido mediante GRPC y se conoce como container runtime interface o CRI. Entre los motores más usados se encuentran containerd y docker. 


## Complementos en Kubernetes

Los complementos en Kubernetes son utilizados para implementar la funcionalidades a nivel de todo el cluster. Existen múltiples complementos disponibles pero por el momento nos vamos a centrar solamente en uno. 

### DNS

El complemento DNS se utiliza CoreDNS. Su función dentro del cluster es gestionar los registros de los servicios asi como la detección y descubrimiento de los mismos. Durante la instalación de Kubernetes viene incorporado este complemento de forma predeterminada. Es estrictamente obligatorio tener esta pieza funcionando para que el cluster funciona correctamente. 


# Enlaces


[1] B. Muschko, CKAD Crash Course. [Online] Available: [CKAD Crash Course](https://github.com/bmuschko/ckad-crash-course) [Accessed: 16-Oct-2022].

[2] B. Muschko, CKA Crash Course. [Online] Available: [CKA Crash Course](https://github.com/bmuschko/cka-crash-course) [Accessed: 16-Oct-2022].

[3] B. Muschko, Certified Kubernetes Application Developer (CKAD) Study Guide. O'Reilly, 2021.

[4] B. Muschko, Certified Kubernetes Application Developer (CKA) Study Guide. O'Reilly, 2022.

[5] J. Domingo, Curso Kubernetes. [Online] Available [Curso Kubernetes](https://www.josedomingo.org/pledin/2022/05/curso-kubernetes/) [Accessed: 16-Oct-2022].

[6] J. Domingo, J.D. Perez, A.M. Coballes, Curso Kubernetes. [Online] Available [Repositorio Curso Kubernetes](https://github.com/iesgn/curso_kubernetes_cep) [Accessed: 16-Oct-2022].

