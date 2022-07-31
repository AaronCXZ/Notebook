- FROM：用来指定要制作的镜像继承自哪个镜像

- MAINTAINER：维护者信息

- RUN：用来执行shell命令

- EXPOSE：将容器的端口暴露出来

- CMD：启动容器时执行的命令，支持三种格式

  1. ```CMD ["exectable", "param1", "param2"]```使用exec执行，推荐
  2. ```CMD command param1 param2```在/bin/sh中执行，提供给需要交换的应用
  3. ```CMD ["param1", "parma2"]```提供给ENTRYPOINT的默认参数

  如果指定了多条CMD指令，只有最后一条被执行。如果用户启动容器指定了运行的命令，会覆盖掉CMD指定的命令

- ENTRYPOINT：支持两种格式

  1. ```ENTRYPOINT ["executable", "param1", "param2"]```
  2. ```ENTRYPOINT command param1 param2```(shell中执行)

  当指定多个时，只有最后一个会被执行

- VOLUME：创建一个可以从本地主机或其他容器挂载的挂载点。

- ENV：指定一个环境变量

- ADD：复制指定的src到容器中的dest，其中src可以时dockerfile所在目录的一个相对路径，也可以是一个URL，还可以是 一个tar文件

- COPY：复制主机的src到容器中的dest，其中src可以时dockerfile所在目录的一个相对路径

