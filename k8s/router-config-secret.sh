echo 'meo' > ./username.txt
echo $ROUTER_PASSWORD > ./password.txt

kubectl create secret generic routercredentials -n router-config \
    --from-file=username=./username.txt \
    --from-file=password=./password.txt

rm ./username.txt ./password.txt