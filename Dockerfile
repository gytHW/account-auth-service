From preugistry.service.ucloud.cn/ibu-common/centos:7.2.1511


COPY ./account_auth /root/account-auth-service/
COPY ./scripts/run.sh /root/account-auth-service/run.sh
COPY ./config/* /root/account-auth-service/config/
COPY ./comic.ttf /root/account-auth-service/comic.ttf

ENTRYPOINT ["/root/account-auth-service/run.sh"]
CMD ["test"]



