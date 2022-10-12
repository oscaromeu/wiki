# Arquitectura básica

## Nodos

k8s es un software que se instala en varios nodos que se gestionan de
forma coordinada, es decir, un clúster de nodos. Aunque es posible en
casos muy peculiares instalar algunos [nodos sobre sistemas
Windows](https://kubernetes.io/docs/setup/production-environment/windows/intro-windows-in-kubernetes/),
la situación normal es que se trate de un cluster de nodos linux. No
es necesario que todos los nodos tengan la misma versión y ni siquiera
que sean la misma distribución, aunque en muchos casos sí lo sea por
simplicidad en el despliegue y mantenimiento.

Los nodos del clúster pueden ser máquinas físicas o virtuales, pero
quizás lo más habitual es que se traten de instancias de nube de
infraestructura, es decir, máquinas virtuales ejecutándose en algún
proveedor de IaaS (AWS, GCP, OpenStack, etc.)

Se distingue entre dos tipos de nodos:

* Los nodos *master*: Son los que ejecutan los servicios principales
de k8s y ordenan a los otros nodos los contenedores que deben
ejecutar. Como el uso del término master es últimamente muy
controvertido en los paises de habla inglesa, se está cambiando su
denominación por *control plane* node.
* Los nodos *worker*: Son los que reciben las órdenes de los
controladores y en los que se ejecutan los contenedores de las
aplicaciones.

## Componentes de un nodo master

* **kube-apiserver** Gestiona la API de k8s
* **etcd** Almacén clave-valor que guarda la configuración del clúster
* **kube-scheduler** Selecciona el nodo donde ejecutar los contenedores
* **kube-controller-manager** Ejecuta los controladores de k8s
* **docker/rkt/containerd/...** Ejecuta los contenedores que sean
  necesarios en el controlador
* **cloud-controller-manager** Ejecuta los controladores que
interactúan con el proveedor de nube:
  * nodos
  * enrutamiento
  * balanceadores
  * volúmenes

## Componentes de un nodo worker

* **kubelet** Controla los Pods asignados a su nodo
* **kube-proxy** Permite la conexión a través de la red
* **docker/rkt/containerd/...** Ejecuta los contenedores
* **supervisord** Monitoriza y controla kubelet y docker

## Complementos (addons)

Los elementos anteriores forman la estructura básica de k8s, pero es
muy habitual que se proporcione funcionalidad adicional a través de
complementos de k8s, que en muchas ocasiones se ejecutan a su vez como
contenedores y son gestionados por el propio Kubernetes. Algunos de
estos complementos son:

* **Cluster DNS** Proporciona registros DNS para los servicios de
  k8s. Normalmente a través de [CoreDNS](https://coredns.io/)
* **Web UI** Interfaz web para el manejo de k8s
* **Container Resource Monitoring** Recoge métricas de forma
centralizada. Múltiples opciones: [prometheus](https://prometheus.io/), [sysdig](https://sysdig.com/)
* **Cluster-level Logging** Almacena y gestiona los logs de los contenedores

## Esquema de nodos y componentes

Se crea un cluster de k8s en los que algunos nodos actúan como master
(normalmente se crea un conjunto impar de nodos master que
proporcionen alta disponibilidad) y el resto actúa como worker en los
que se ejecutan los contenedores de las aplicaciones. Los nodos se
comunican entre sí a través de una red que proporciona la capa de
infraestructura y se crea una red para la comunicación de los
contenedores, que suele ser una red de tipo
[overlay](https://en.wikipedia.org/wiki/Overlay_network).

<img src="https://github.com/iesgn/curso_kubernetes_cep/raw/main/modulo1/img/arquitectura.png" alt="arquitectura" />

## Vídeo

[https://www.youtube.com/watch?v=mF-OGBbA57k](https://www.youtube.com/watch?v=mF-OGBbA57k)

# Docker

Docker es una empresa ([Docker Inc.](https://www.docker.com/)) que
desarrolla un software con el mismo nombre, de forma más concreta el software denominado ([docker
engine](https://www.docker.com/products/container-runtime)), que ha
supuesto una revolución en el desarrollo de software, muy ligado al
uso de contenedores de aplicaciones, a las aplicaciones web y al
desarrollo ágil.

Docker permite gestionar contenedores a alto nivel, proporcionando
todas las capas y funcionalidad adicional y, lo más importante de todo,
es que proporciona un nuevo paradigma en la forma de distribuir las
aplicaciones, ya que se crean imágenes en contenedores que se
distribuyen, de manera que el contenedor que se ha desarrollado es
idéntico al que se utiliza en producción y deja de instalarse la
aplicación de forma tradicional.

## Componentes de docker

Docker engine tiene los componentes que a *grosso modo* se presentan a
continuación:

<img src="https://github.com/iesgn/curso_kubernetes_cep/raw/main/modulo1/img/docker.png" alt="docker" />

En la imagen se han destacado los componentes que son relevantes desde
el punto de vista de este curso, ya que como veremos más adelante,
docker podría ser un componente esencial de Kubernetes, pero realmente
no lo es completo, solo containerd y los elementos que éste
proporciona lo son, ya que k8s utiliza su propia API, su propia línea
de comandos y gestiona el almacenamiento y las redes de forma
independiente a docker.

## Evolución del proyecto docker

Docker tuvo un enorme éxito y una gran repercusión, pero la empresa
que lo desarrolla siempre se ha movido en el dilema de cómo sacar
rendimiento económico a su software, que al ser desarrollado bajo
licencia libre, no proporciona beneficio como tal. Este dilema se ha
tratado de resolver con modificaciones en la licencia o con doble
licenciamiento (docker CE y docker EE en estos momentos), pero esto a
su vez ha propiciado que otras empresas desarrollasen alternativas a
docker para no depender en el futuro de una empresa sin un modelo de
negocio claro y ante posibles modificaciones de la licencia libre de
docker.

Los cambios más significativos que han ocurrido en docker se enumeran
a continuación:

* [Moby](https://github.com/moby/moby) Docker engine se desarrolla
ahora como proyecto de software libre independiente de Docker Inc. denominándose Moby. De este proyecto se surten las distribuciones de
linux para desarrollar los paquetes docker.io

* [Docker Engine](https://www.docker.com/products/container-runtime)
Versión desarrollada por Docker Inc.

* [runC](https://github.com/opencontainers/runc) Componente que
ejecuta los contenedores a bajo nivel. Actualmente desarrollado por
OCI

* [containerd](https://github.com/containerd/containerd) Componente
que ejecuta los contenedores e interactúa con las
imágenes. Actualmente desarrollado por la CNCF.

## Limitaciones de docker (docker engine)

Docker (docker engine) gestiona completamente la ejecución de un
contenedor en un determinado nodo a partir de una imagen, pero no
proporciona toda la funcionalidad que necesitamos para ejecutar
aplicaciones en entornos en producción.

Existen diferentes preguntas
que nos podemos hacer acerca de esto :

* ¿Qué hacemos con los cambios entre versiones?
* ¿Cómo hacemos los cambios en producción?
* ¿Cómo se balancea la carga entre múltiples contenedores iguales?
* ¿Cómo se conectan contenedores que se ejecuten en diferentes
demonios de docker?
* ¿Se puede hacer una actualización de una aplicación sin
interrupción?
* ¿Se puede variar a demanda el número de réplicas de un determinado
contenedor?
* ¿Es posible mover la carga entre diferentes nodos?

Las respuestas a estas preguntas no pueden venir de docker engine, ya
que no es un software desarrollado para eso, tiene que venir de algún software
que pueda utilizar docker o parte de él y que sea capaz de comunicar
múltiples nodos para proporcionar de forma coordinada estas
funcionalidades. Ese software se conoce de forma genérica como
**orquestador de contenedores**.

## Vídeo

[https://www.youtube.com/watch?v=UdPsknw30OE](https://www.youtube.com/watch?v=UdPsknw30OE)


# El proyecto Kubernetes

El proyecto Kubernetes lo inicia Google en 2014 como un software
(libre) para orquestar contenedores. En aquel momento había varios
proyectos de software que querían extender las posibilidades del uso
de contenedores de aplicaciones tipo docker a entornos en producción,
lo que de forma genérica se conoce como orquestadores de
contenedores. A diferencia del resto, Kubernetes no es un proyecto que
se desarrolla desde cero, sino que aprovecha todo el conocimiento que
tenía Google con el uso de la herramienta interna
[Borg](https://kubernetes.io/blog/2015/04/borg-predecessor-to-kubernetes/),
de manera que cuando se hace pública la primera versión de Kubernetes,
ya era un software con muchas funcionalidades.

Un proyecto se convierte en software libre cuando utiliza una
[licencia libre](https://opensource.org/licenses/category), pero otro
aspecto importante es la gobernanza del proyecto, es decir, si el
desarrollo es abierto o no, si las decisiones sobre las nuevas
funcionalidades las toma una empresa o se consensúan, etc. Si un
proyecto de software libre lo inicia una única empresa, siempre existe
la desconfianza de que ese proyecto vaya a ir encaminado a beneficiar
a esa empresa. En este caso, la empresa en cuestión era un gigante
como Google, por lo que aunque el proyecto era muy interesante,
existía cierto recelo de gran parte del sector inicialmente. Para
conseguir que una parte importante del sector se sumase al proyecto,
Google tomó la decisión de desvincularse del mismo y ceder el control
a la [Cloud Native Compute Foundation (CNCF)](https://www.cncf.io/),
por lo que Kubernetes es un proyecto de software libre de fundación,
en el que se admiten contribuciones de forma abierta y donde las reglas de
la gobernanza recaen sobre los [miembros de la
fundación](https://www.cncf.io/about/members/), normalmente un
conjunto amplio de grandes empresas del sector. Es decir, aunque hoy
en día hay quien habla de Kubernetes como el software de orquestación
de contenedores de Google, esto es un error, es un proyecto que
gestiona desde hace años la CNCF, a la que ni siquiera pertenece
Google.

## ¿Qué es Kubernetes?

Kubernetes es un software pensado para gestionar completamente el
despliegue de aplicaciones sobre contenedores, realizando este
despliegue de forma completamente automática y poniendo un gran
énfasis en la escalabilidad de la aplicación, así como el control
total del ciclo de vida. Por destacar algunos de los puntos más
importantes de Kubernetes, podríamos decir:

* Despliega aplicaciones rápidamente
* Escala las aplicaciones al vuelo
* Integra cambios sin interrupciones
* Permite limitar los recursos a utilizar

Kubernetes está centrado en la puesta en **producción** de
contenedores y por su gestión es indicada para administradores de
sistemas y personal de equipos de operaciones. Por otra parte, afecta también
a los desarrolladores, ya que las aplicaciones deben adaptarse para
poder desplegarse en Kubernetes.

## Características principales

<img src="https://github.com/iesgn/curso_kubernetes_cep/raw/main/modulo1/img/logo.png" alt="k8s-logo" width="150" />

Kubernetes surge como un software para desplegar aplicaciones sobre
contenedores que utilicen infraestructura en nube (pública, privada o
híbrida). Aunque puede desplegarse también en entornos más
tradicionales como servidores físicos o virtuales, no es su "entorno
natural".

Kubernetes es extensible, por lo que cuenta con gran cantidad de
módulos, plugins, etc.

El nombre del proyecto proviene de una palabra de griego antiguo que
significa timonel y habitualmente se escribe de forma abreviada como
k8s.

## Características del software

Kubernetes está desarrollado en el lenguaje [Go](https://golang.org/)
como diversas aplicaciones de este sector. La primera versión de
Kubernetes se publicó el 7 de junio de 2014, aunque la más antigua
disponible en el repositorio es la
[v0.2](https://github.com/kubernetes/kubernetes/releases/tag/v0.2), de
septiembre de 2014.

La licencia utilizada en Kubernetes es la [Apache License
v2.0](https://www.apache.org/licenses/LICENSE-2.0.html), licencia de
software libre permisiva, muy utilizada últimamente en proyectos de
fundación en los que están involucrados empresas, ya que no se trata
de una licencia copyleft, que no permitiría su inclusión en software
que no sea libre, mientras que la licencia Apache sí lo permite en
determinadas circunstancias.

El código de Kubernetes se gestiona a través de
[Github](https://github.com/kubernetes/kubernetes) en cuyo repositorio
se puede ver la gran cantidad de código desarrollado en estos años
(más de 100000 "commits") y las miles de personas que han participado
en mayor o menor medida. La última versión de Kubernetes en el momento
de escribir esta documentación es la 1.23 y el proyecto actualmente
está publicando dos o tres versiones nuevas cada año.

En cualquier caso la versión de Kubernetes no es algo esencial para
los contenidos de este curso, porque se van a tratar los elementos
básicos, que ya están muy establecidos y, salvo algún detalle menor,
se puede realizar este curso al completo con una versión de Kubernetes
diferente a la utilizada para la documentación.

## El ecosistema

De entre todas las opciones de orquestadores de contenedores
disponibles, hoy se considera que la opción preferida en la mayor
parte de los casos es k8s y se ha desarrollado un enorme ecosistema de
aplicaciones alrededor que proporcionan algunas funcionalidades que no
tiene k8s o que de alguna forma utiliza o se pueden integrar de
diferente forma con k8s. Este ecosistema de aplicaciones está
actualmente en plena "ebullición" y es posible que en unos años
algunos de esos proyectos se estabilicen y otros desaparezcan, ya que
en muchos casos solapan unos con otros.

[https://landscape.cncf.io/](https://landscape.cncf.io/)

## Vídeo

[https://www.youtube.com/watch?v=MtA74Hc4FAo](https://www.youtube.com/watch?v=MtA74Hc4FAo)