FROM golang:1.21.1-bullseye

ARG USERNAME=noda
ARG GROUP=noda
ARG UID=1001
ARG GID=1001

RUN groupadd -g ${GID} ${GROUP} && \ 
    adduser --disabled-password --gecos '' --gid ${GID} --uid ${UID} ${USERNAME} && \
    mkdir -p /etc/sudoers.d && echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/${USERNAME} && \
    chmod 0440 /etc/sudoers.d/${USERNAME}

USER $USERNAME
WORKDIR /home/${USERNAME}/app
