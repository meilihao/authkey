# init-env
```bash
# --- set bash
# cat << "EOF" > init-env.sh
#!/usr/bin/env bash
if [ "`(cat ~/.bashrc |grep \"DISABLE_UPDATE_PROMPT\") 2> /dev/null`" = "" ]
then
    # disable oh-my-bash show update prompt where ssh connect
    echo "DISABLE_UPDATE_PROMPT=true" >> ~/.bashrc
    echo "DISABLE_AUTO_UPDATE=true" >> ~/.bashrc
fi
EOF
# chmod +x init-env.sh
# ./init-env.sh
# --- set hostname
# hostnamectl set-hostname <newhostname> # 更新/etc/hosts文件
# vim /etc/hosts
```