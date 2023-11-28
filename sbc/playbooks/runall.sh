#/bin/bash
for f in ./*.yml
do
ansible-playbook $f
done
