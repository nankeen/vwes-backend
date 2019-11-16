FROM alpine:3.7
COPY vwes-backend /app/vwes-backend
ENV PORT 80
EXPOSE 80
WORKDIR /app
ENTRYPOINT ["/app/vwes-backend"]
