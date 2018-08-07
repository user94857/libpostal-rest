# Libpostal REST

## API Example

Replace `<host>` with your host ip address

### Parser
### plain text in post body:
`curl -X POST -d '100 main st buffalo ny' <host>:8080/parser`

** Response **
```
[
  {
    "label": "house_number",
    "value": "100"
  },
  {
    "label": "road",
    "value": "main st"
  },
  {
    "label": "city",
    "value": "buffalo"
  },
  {
    "label": "state",
    "value": "ny"
  }
]
```

### Added options to this api which will be passed to the libpostal library
`curl -X POST -d '100 main st buffalo ny' <host>:8080/parser?country=us`

### Or
`curl -X POST -d '100 main st buffalo ny' <host>:8080/parser?language=en`

### Or
`curl -X POST -d '100 main st buffalo ny' <host>:8080/parser?language=en&country=us`

## NEW GET method
## Due to the way mux handles encoded/decoded urls, I'm unable to use an encoded address in the url
## This means the request will not parse properly if you don't 'clean' the address
## Replace all non alpha numeric letters including white spaces with a comma. Except for - . _ / = that i've tested so far. Those characters are Ok. It seems to work just as well
`curl -G <host>:8080/parser?address=100,main,st,buffalo,ny`

In your browser:
`http://<host>:8080/parser?address=100,main,st,buffalo,ny`

### Or
`curl -G <host>:8080/parser?address=100,main,st,buffalo,ny&language=en`

### Or
`curl -G <host>:8080/parser?address=100,main,st,buffalo,ny&country=us`

### Or
`curl -G <host>:8080/parser?address=100,main,st,buffalo,ny&language=en&country=us`


### Expand
### no formatting just plaintext address:
`curl -X POST -d '100 main st buffalo ny' <host>:8080/expand`

** Response **
```
[
  "100 main saint buffalo new york",
  "100 main saint buffalo ny",
  "100 main street buffalo new york",
  "100 main street buffalo ny"
]
```
