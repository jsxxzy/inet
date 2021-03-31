
export dr_user=用户
export dr_password=密码

cd /tmp
mkdir app
cd app
mv ../inetd .

echo "添加执行权限"
chmod u+x inetd

echo "后台运行"
nohup ./inetd 1 >/dev/null 2>&1 </dev/null &

echo "删除日志文件"
rm -rf nohup.out
