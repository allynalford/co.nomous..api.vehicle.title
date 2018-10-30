
FROM golang:1.7.4-onbuild


RUN go get github.com/gin-gonic/gin

RUN go get github.com/go-sql-driver/mysql

RUN go get gopkg.in/mgo.v2



# Document that the service listens on port 8080.
#EXPOSE 5000



#RUN go build ./application.go

#CMD ["cd","./app"]

#CMD ["go", "run", "application.go"]



#RUN go build  application.go

COPY ./com.plydge.api.vehicle.title /app/com.plydge.api.vehicle.title
RUN chmod +x /app/com.plydge.api.vehicle.title

ENV PORT 8080
EXPOSE 8080

ENTRYPOINT /app/com.plydge.api.vehicle.title


