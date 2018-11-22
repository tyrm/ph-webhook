FROM scratch
LABEL maintainer="tyr@pettingzoo.co"

EXPOSE 8080

ENV DB_ENGINE  postgresql://postgres:5432/api

ADD ph-webhook /
CMD ["/ph-webhook"]
