

# Go API Container
FROM golang:1.13.9-alpine AS builder
ADD . /app   
WORKDIR /app/server
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main .

# React application Container
FROM node:alpine AS node_builder
COPY --from=builder /app/webapp ./
COPY --from=builder /app/.env ./
RUN npm install
RUN npm audit fix
RUN npm run build #Final stage build, this will be the container

# that we will deploy to production
FROM alpine:3.11.5
RUN apk --no-cache add ca-certificates
COPY --from=builder /main ./
COPY --from=node_builder /build ./web
RUN chmod +x ./main

#EXPOSE API PORT

EXPOSE 7777
CMD ["./main"]


#COPY --from=builder /app/.env .       
# 
# 
# 
# 
# /api
#   /webapp
# 