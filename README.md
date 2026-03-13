# Zelda JSON API (Go)

## Descripción

Esta es una API REST que permite consultar y gestionar información sobre juegos de Zelda.

Los datos se almacenan en un archivo JSON.

El servidor corre en el puerto correspondiente a mi número de carnet: 24979

# Estructura del proyecto

go-http
│
├── main.go
├── data
│   └── items.json
└── README.md

# Endpoints


## Obtener un juego por query parameter

GET

```
/api/items?id=3
```

Devuelve un juego específico según su ID.

Ejemplo:

```
http://localhost:24979/api/items?id=3
```


## Obtener un juego por path parameter

GET

```
/api/items/3
```

Devuelve un juego específico según su ID.


## Crear un nuevo juego

POST

```
/api/items/create
```

Body JSON de ejemplo:

```json
{
 "title": "Oracle of Seasons",
 "console": "Game Boy Color",
 "year": 2001,
 "developer": "Nintendo",
 "genre": "Action Adventure"
}
```


## Eliminar un juego

DELETE

```
/api/items/delete/{id}
```

Ejemplo:

```
/api/items/delete/5
```

# Filtros disponibles

La API permite filtrar resultados utilizando query parameters.

Ejemplos:

Filtrar por consola:

```
/api/items?console=Switch
```

Filtrar por año:

```
/api/items?year=1998
```

Filtrar por consola y año:

```
/api/items?console=Switch&year=2017
```

---

# Manejo de errores

La API devuelve códigos HTTP apropiados en caso de error

# Cómo ejecutar el proyecto

1. Clonar el repositorio

```
git clone <repo>
```

2. Entrar al proyecto

```
cd go-http
```

3. Ejecutar el servidor

```
go run main.go
```

4. Abrir en navegador o Postman

```
http://localhost:24979/api/items
```

# Evidencia de pruebas

Las pruebas de los endpoints fueron realizadas utilizando **Postman** incluyendo:

* GET de todos los elementos
* GET con query parameters
* POST exitoso
* eliminación con DELETE
* manejo de errores
