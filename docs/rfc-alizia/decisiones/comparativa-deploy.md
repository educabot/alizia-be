# Cloud Run vs Cloud Functions — Comparativa para decisión del equipo

## ¿Qué es cada uno?

| | Cloud Functions | Cloud Functions agrupadas | Cloud Run |
|---|---|---|---|
| **Qué subís** | 1 función Go | 1 función por módulo | 1 container Docker |
| **Qué corre** | 1 endpoint por función | Varios endpoints por función | Todos los endpoints en 1 proceso |
| **Ejemplo** | 55+ funciones (como tich-cronos hoy) | 5 funciones (1 por módulo) | 1 container |
| **Analogía** | 55 microaplicaciones | 5 miniservices | 1 servidor normal en la nube |

---

## Comparativa detallada

| Aspecto | Cloud Functions (1:1) | Cloud Functions agrupadas | Cloud Run |
|---|---|---|---|
| **Unidad de deploy** | 1 función = 1 endpoint | 1 función = 1 módulo (N endpoints) | 1 container = todos los endpoints |
| **Cantidad de deploys** | 55+ (como tich-cronos) | 5 (admin, coordination, teaching, resources, ai) | 1 solo |
| **Cold start** | Si, en cada función | Si, pero solo 5 funciones | Configurable (min instances = 1) |
| **Latencia cold start** | ~500ms - 2s (Go) | ~500ms - 2s (Go) | ~200ms - 500ms (container buildeado) |
| **Escalado** | Por endpoint individual | Por módulo | Por container (réplicas) |
| **Dockerfile** | No | No | Sí (container estándar) |
| **Portabilidad** | Solo GCP | Solo GCP | Cualquier cloud, VPS, Docker |
| **Complejidad operativa** | Alta (55+ registros) | Media (5 registros) | Baja (1 deploy, 1 config) |
| **Debugging local** | Emulador o Go local | Emulador o Go local | `docker run` o `go run` |
| **Logs** | 55+ streams | 5 streams (1 por módulo) | 1 stream unificado |
| **CI/CD** | Deploy selectivo o todo | Deploy por módulo | 1 pipeline, 1 build, 1 deploy |
| **Rollback** | Por función individual | Por módulo | Todo junto |
| **Max request timeout** | 9 min (Gen2) | 9 min (Gen2) | 60 min |
| **Max instances** | 1000 por función | 1000 por función | 1000 containers |
| **Min instances** | 0 o 1+ por función | 0 o 1+ por función | 0 o 1+ por container |
| **Conexiones DB** | 55+ pools (saturan PostgreSQL) | 5 pools (manejable) | 1 pool compartido |
| **Estado en memoria** | No (stateless) | Limitado (dentro del módulo) | Sí (cache, connection pool) |

---

## Costos

| Escenario | Cloud Functions (1:1) | Cloud Functions agrupadas | Cloud Run |
|---|---|---|---|
| **Tráfico bajo** (< 2M req/mes) | Más barato. Free tier generoso | Más barato. Free tier generoso | Más caro si min instances = 1 |
| **Tráfico medio** (2M - 10M req/mes) | Similar | Similar | Similar |
| **Tráfico alto** (> 10M req/mes) | Más caro (overhead por invocación) | Más caro (overhead por invocación) | Más barato (container absorbe) |
| **Idle** (sin tráfico) | $0 | $0 | $0 o ~$5-15/mes (min instances) |
| **Free tier** | 2M invocaciones, 400K GB-sec | 2M invocaciones, 400K GB-sec | 2M requests, 360K vCPU-sec |

---

## Pros y Contras

### Cloud Functions

**Pros:**
1. Free tier generoso para bajo tráfico
2. Escalado granular por endpoint (si /generate necesita más, solo esa escala)
3. No necesitás Dockerfile
4. Aislamiento: si un endpoint crashea, los otros siguen
5. Ya lo usa tich-cronos (experiencia del equipo)
6. Pay-per-use real (0 tráfico = $0)

**Contras:**
1. Cold starts en cada función (mala UX en primera request)
2. 55+ funciones = 55+ deploys = 55+ configs = mantenimiento tedioso
3. Cada función abre su propio pool de DB (puede saturar PostgreSQL)
4. No portable (solo GCP)
5. Logs fragmentados (1 stream por función)
6. CI/CD más complejo (deploy selectivo o deploy de todo)
7. Debugging local requiere emulador
8. Registrar cada función nueva es manual y propenso a errores

### Cloud Functions agrupadas por módulo (opción intermedia)

**Pros:**
1. Solo 5 funciones en vez de 55+ (admin, coordination, teaching, resources, ai)
2. Free tier de Cloud Functions (bajo costo)
3. Escalado por módulo (si AI necesita más, solo esa función escala)
4. Misma experiencia del equipo (ya conocen Cloud Functions)
5. 5 pools de DB en vez de 55+ (no satura PostgreSQL)
6. Si un módulo crashea, los otros siguen
7. Deploy por módulo (cambio en teaching no redeploya admin)
8. Mejora incremental sobre lo que ya tienen

**Contras:**
1. Cold starts siguen existiendo (5 funciones que arrancan)
2. Solo GCP (no portable)
3. 5 registros de funciones para mantener (poco pero no cero)
4. Cada función sigue teniendo su propio pool de DB
5. Timeout máximo 9 min (AI generation larga puede ser problema)

### Cloud Run

**Pros:**
1. 1 deploy = todos los endpoints
2. Cold start mínimo con min instances = 1
3. Container estándar → portable a AWS, Azure, o VPS si algún día migran
4. 1 pool de conexiones DB compartido (eficiente)
5. Cache in-memory posible (entre requests)
6. Logs unificados (1 stream)
7. CI/CD simple: build, push, deploy
8. Debugging local: `docker run` y listo
9. Misma experiencia dev y prod

**Contras:**
1. Con min instances = 1 pagás idle (~$5-15/mes por servicio)
2. Escalado es todo-o-nada (no por endpoint individual)
3. Necesitás Dockerfile (aunque es simple)
4. Si el container crashea, caen TODOS los endpoints
5. Rollback es todo junto (no por endpoint)

---

## Para nuestro caso específico (Alizia)

| Factor | Cloud Functions (1:1) | Cloud Functions agrupadas | Cloud Run | Observación |
|---|---|---|---|---|
| **26+ tablas, JOINs complejos** | 55+ pools saturan DB | 5 pools manejable | 1 pool compartido | Agrupadas o Cloud Run mejor para DB-heavy |
| **AI generation (10-30s)** | Función ocupada esperando | Función ocupada esperando | Container maneja concurrencia | Cloud Run mejor para long-running |
| **4+ devs** | 55+ funciones para coordinar | 5 funciones claras | 1 deploy coordinado | Agrupadas buen balance |
| **Experiencia del equipo** | Ya lo conocen (tich-cronos) | Mismo approach mejorado | Curva mínima (Docker) | Agrupadas = mejora incremental |
| **Presupuesto limitado** | Free tier generoso | Free tier generoso | Min instances tiene costo | Cloud Functions si ajusta el presupuesto |
| **Portabilidad futura** | Solo GCP | Solo GCP | Cualquier cloud/VPS | Cloud Run si quieren flexibilidad |
| **Simplicidad operativa** | Baja (55+ configs) | Media (5 configs) | Alta (1 config) | Cloud Run más simple |

---

## Opción intermedia: Cloud Functions con pocas funciones

En vez de 55+ funciones (1 por endpoint), agrupar por módulo:

```
b2b-alizia-http-admin          → Todos los endpoints de /admin/*
b2b-alizia-http-coordination   → Todos los endpoints de /coordination-documents/*
b2b-alizia-http-teaching       → Todos los endpoints de /lesson-plans/*
b2b-alizia-http-resources      → Todos los endpoints de /resources/*
b2b-alizia-http-ai             → Todos los endpoints de /chat, /generate
```

5 funciones en vez de 55+. Reduce la complejidad operativa de Cloud Functions manteniendo el escalado granular.

---

## Resumen para la decisión

| Si el equipo prioriza... | Elegir |
|---|---|
| Experiencia existente + bajo costo + no cambiar nada | Cloud Functions (1:1) |
| Mejorar lo actual sin cambio radical + bajo costo | Cloud Functions agrupadas por módulo |
| Simplicidad operativa + portabilidad + long-running | Cloud Run |
