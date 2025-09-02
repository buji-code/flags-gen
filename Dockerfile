FROM scratch

COPY flags-gen /usr/local/bin/flags-gen

ENTRYPOINT ["/usr/local/bin/flags-gen"]