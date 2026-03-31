# HU-5.6: Bitácora de cotejo

> Como docente, necesito registrar cómo fue la clase después de dictarla, para que el sistema aprenda y mejore las propuestas futuras.

**Fase:** Fase 5
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] El docente puede registrar la bitácora después de una clase planificada
- [ ] Soporta entrada por audio (grabación libre) y por texto
- [ ] El audio se transcribe automáticamente usando el servicio de IA
- [ ] La bitácora incluye: qué se cumplió, qué no, observaciones sobre alumnos, ajustes necesarios
- [ ] Alizia puede hacer preguntas de seguimiento si detecta información faltante
- [ ] La bitácora queda vinculada al lesson plan de la clase
- [ ] El coordinador puede ver las bitácoras de su equipo

## Concepto

La bitácora es el cierre del ciclo planificación → ejecución → reflexión. El docente graba un audio libre contando cómo fue la clase — sin formato rígido, porque la adopción depende de que sea natural y rápido.

### Flujo

```
Docente termina la clase
  → Abre la app en el celular
  → Graba audio: "Hoy hicimos la actividad de grupo pero..."
  → El sistema transcribe y estructura
  → Alizia pregunta: "¿Llegaste a hacer el cierre?"
  → Docente responde
  → Bitácora guardada y vinculada al plan
```

### Datos que se recolectan

- ¿Se cumplió el objetivo? (sí/parcial/no)
- ¿Se completaron todas las actividades?
- Observaciones sobre alumnos específicos (seguimiento)
- Ajustes sugeridos por el docente para futuras clases
- Tiempo real vs. planificado

### Consideraciones

- **Offline:** ¿Qué pasa si el docente no tiene internet? Escuelas rurales argentinas pueden tener conectividad limitada. ¿Se sube después? ¿Funciona offline? → pendiente definición
- **Privacidad:** Las observaciones sobre alumnos son sensibles y requieren política de acceso clara
- **Adopción:** La clave es que sea rápido y natural. Si toma más de 2-3 minutos, los docentes no lo usarán

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 5.6.1 | [Modelo de datos bitácora](./tareas/T-5.6.1-modelo-datos-bitacora.md) | entities, migracion | ⬜ |
| 5.6.2 | [Transcripción de audio](./tareas/T-5.6.2-transcripcion-audio.md) | providers/audio_transcriber.go | ⬜ |
| 5.6.3 | [Procesamiento estructurado y seguimiento](./tareas/T-5.6.3-procesamiento-estructurado.md) | usecases/ai/ | ⬜ |
| 5.6.4 | [Endpoints CRUD de bitácora](./tareas/T-5.6.4-endpoints-bitacora.md) | handlers/ | ⬜ |
| 5.6.5 | [Tests](./tareas/T-5.6.5-tests.md) | tests/ | ⬜ |

## Dependencias

- [HU-5.5: Estados](../HU-5.5-estados-planificacion/HU-5.5-estados-planificacion.md) — Solo clases published/dictadas
- [Épica 6: Asistente IA](../../06-asistente-ia/06-asistente-ia.md) — Transcripción y procesamiento de audio

## Notas

> **Nota:** La "recoleccion activa de datos" que producto menciona (Alizia pregunta por informacion faltante o alumnos con seguimiento) se implementa dentro del chat de la bitacora. No es un flujo autonomo separado — Alizia hace estas preguntas como parte de la conversacion cuando detecta datos faltantes en el contexto.
