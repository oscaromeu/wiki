---
id: observabilidad
title: Observabilidad
sidebar_label: Observabilidad
sidebar_position: 11
draft: true
---

## Realizar checks para verificar el correcto funcionamiento de las aplicaciones

Una vez hemos desplegado nuestras aplicaciones en el cluster nos interesa saber si estas se estan ejecutando de manera correcta ya sea después de una hora, una semana, un mes o cuando sea. En esta sección vamos a dar unas pinceladas de observabilidad.

## Health Probing 

Incluso con las mejores prácticas de ingenieria de test automatizado es practicamente imposible que detectemos todos los bugs antes de ponerlo en producción. Esto es especialmente cierto para situaciones en las que el software ha estado ejecutandose durante largos periodos de tiempo en el cluster, no es poco frecuente que nos encontremos ante situaciones de memory leaks, deadlocks, bucles infintos y situaciones similares que suceden cuando la aplicación ya lleva corriendo durante un tiempo en producción. 

Monitorizar de manera adecuada los sistemas es un primer paso para prevenir las situaciones descritas pero aún así es necesario coger el toro por los cuernos para remediar dichas situaciones.

Kubernetes nos ofrece un concepto llamado "Health Probing" que nos permite de manera periodica y automatizada comprobar el estado de salud de nuestras aplicaciones. Podemos configurar nuestros contenedores para que ejecuten procesos livianos que comprueben diferentes casuisticas.

### Readiness probes

Una vez nuestra aplicación ya esta puesta en marcha puede necesitar algunas configuraciones adicionales. Por ejemplo, conectarse a una base de datos y realizar algún tipo de configuración. Este tipo de sonda comprueba si nuestra aplicación esta lista para recibir tráfico. 


![readiness probes](./11/img/01-readiness-probes.png#center)

#### Definir una Readiness probe

```yaml
apiVersion: v1
kind: Pod
metadata:
 name: web-app
spec:
 containers:
 - name: web-app
 image: eshop:4.6.3
 readinessProbe:
 httpGet:
 path: /
 port: 8080
 initialDelaySeconds: 5
 periodSeconds: 2
```

### Liveness probes

Una vez nuestra aplicación ya esta puesta en marcha queremos que siga funcionando sin nigún tipo de error. Esta sonda realiza comprobaciones de manera periodica para asegurar tal condición. 

![liveness probes](./11/img/02-livenessprobes.png#center)

#### Definir una Liveness probe

```
apiVersion: v1
kind: Pod
metadata:
 name: web-app
spec:
 containers:
 - name: web-app
 image: eshop:4.6.3
 livenessProbe:
 exec:
 command:
 - cat
 - /tmp/healthy
 initialDelaySeconds: 10
 periodSeconds: 5
```


### Startup probes

Las aplicaciones legacy pueden requerir algún tiempo adicional antes de que se empieze a realizar algún tipo de acción en ellos. Esta sonda sirve para dar un tiempo extra de arranque y evitar que otras sondas (liveness-probe) empiezen a realizar sus correspondientes comprobaciones. Esta sonda elimina el contenedor si este no arranca en el tiempo marcado por la sonda de startup. 

![startup probes](./11/img/03-startup-probes.png#center)

#### Definir una startup probe

```yaml
apiVersion: v1
kind: Pod
metadata:
 name: startup-pod
spec:
 containers:
 - image: httpd:2.4.46
 name: http-server
 startupProbe:
 tcpSocket:
 port: 80
 initialDelaySeconds: 3
 periodSeconds: 15

```


Cada tipo de sonda nos ofrece tres métodos diferentes para poder realizar las comprobaciones

![health verification methods](./11/img/health-verification-methods.png#center)

cada sonda tiene una serie de atributos que nos permiten ajustar el comportamiento deseado

![health verification checks](./11/img/health-verification-methods.png#center)

### Debugging en Kubernetes

### Monitoring


### Optimizar recursos de CPU y memoria

