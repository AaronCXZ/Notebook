#### 一个简单的playbook

```yaml
---
- name: Configure webserver with nginx
  hosts: webservers
  sudo: True
  vars:
    key_file: /etc/nginx/ssl/nginx.key
    cert_file: /etc/nginx/ssl/nginx.crt
    conf_file: /etc/nginx/sites-available/default
    server_name: localhost
  tasks:
  - name: install nginx
    age: name=nginx update_cache=yes cache_valid_time=3600
    
  - name: create directories for ssl certificates
    file: path=/etc/nginx/ssl state=directory
    
  - name: copy TLS key
    copy: src=files/nginx.key dest={{ key_file }} owner=root mode=0600
    notify: restart nginx
    
  - name: copy TLS certificate
    copy: src=files/nginx.crt dest={{ cert_file }}
    notify: restart nginx

  - name: copy nginx config file
    template: src=templates/nginx.conf.j2 dest={{ conf_file }}
    notify: restart nginx

  - name: enable configuration
    file: >
      dest=/etc/nginx/sites-enabled/default
      src={{ conf_file }}
      stat=link
    notify: restart nginx

  - name: copy index.html
    template: src=templates/index.html.j2 dest=/usr/share/nginx/html/index/html mode=0644

  handlers:
  - name: restart nginx
    service: name=nginx state=restarted
```

Gathering Facts 的作用：ansible开始运行playbook的时候第一件事就是连接到服务器上收集各种信息，包括：操作系统、主机名、所有网络即可的IP地址和MAC地址等，这样就可以在之后的playbook中使用这些信息了，如果不需要这些信息可以关闭以节省一些时间。

handlers：只有在被task通知时才会运行，且只会在所有任务都执行完之后执行，而且即使被通知了多次，也只会执行一次。handler执行的顺序是按play的顺序，不是被通知的顺序。一般用来重启服务。

ansible变量的优先级（由高到低）：

1. ansible-playbook -e var=value
2. 这个优先级列表中没有提到的其他方法
3. 通过inventory文件或者YAML文件定义的主机变量或群组变量
4. Fact
5. 在role的defaults/mail.yml文件中

 