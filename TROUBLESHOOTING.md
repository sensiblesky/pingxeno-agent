# Troubleshooting SCP Connection Issues

## Common Issues and Solutions

### 1. Test SSH Connection First

Before using SCP, test if SSH works:

```bash
ssh cyb4x@192.168.1.100
```

If SSH works but SCP doesn't, try the solutions below.

### 2. Check SSH Service on Remote Server

On the remote server, verify SSH is running:

```bash
# On the Linux server
sudo systemctl status ssh
# or
sudo systemctl status sshd
```

If not running:
```bash
sudo systemctl start ssh
sudo systemctl enable ssh
```

### 3. Check Firewall

On the remote server, ensure port 22 is open:

```bash
# Ubuntu/Debian (ufw)
sudo ufw allow 22/tcp
sudo ufw status

# CentOS/RHEL (firewalld)
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --reload
```

### 4. Try with Verbose Output

Get more details about the connection failure:

```bash
scp -v pingxeno-agent-linux-amd64 cyb4x@192.168.1.100:/tmp
```

### 5. Try Different Port

If SSH is on a non-standard port:

```bash
scp -P 2222 pingxeno-agent-linux-amd64 cyb4x@192.168.1.100:/tmp
```

### 6. Alternative Transfer Methods

#### Option A: Using SFTP

```bash
sftp cyb4x@192.168.1.100
# Once connected:
put pingxeno-agent-linux-amd64 /tmp/
exit
```

#### Option B: Using rsync

```bash
rsync -avz -e ssh pingxeno-agent-linux-amd64 cyb4x@192.168.1.100:/tmp/
```

#### Option C: Using HTTP Server (Python)

On your local machine:
```bash
cd /Users/denicsann/Desktop/projects/pingxeno/agent
python3 -m http.server 8001
```

On the remote server:
```bash
wget http://192.168.1.X:8001/pingxeno-agent-linux-amd64 -O /tmp/pingxeno-agent
# Replace 192.168.1.X with your local machine's IP
```

#### Option D: Using USB Drive

1. Copy the binary to a USB drive
2. Plug it into the server
3. Mount and copy:
```bash
sudo mount /dev/sdb1 /mnt
cp /mnt/pingxeno-agent-linux-amd64 /tmp/pingxeno-agent
```

#### Option E: Using Base64 Encoding (for small files)

On local machine:
```bash
base64 pingxeno-agent-linux-amd64 > agent_base64.txt
```

Then on remote server, paste the content and:
```bash
base64 -d agent_base64.txt > /tmp/pingxeno-agent
chmod +x /tmp/pingxeno-agent
```

### 7. Check Network Connectivity

```bash
# Ping the server
ping 192.168.1.100

# Check if port 22 is open
nc -zv 192.168.1.100 22
# or
telnet 192.168.1.100 22
```

### 8. SSH Configuration Issues

Check your SSH config:

```bash
cat ~/.ssh/config
```

Try with explicit options:

```bash
scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null pingxeno-agent-linux-amd64 cyb4x@192.168.1.100:/tmp
```

### 9. Check Permissions

Ensure you have write permissions on the remote server:

```bash
# On remote server
ls -la /tmp
# Should show write permissions
```

### 10. Try Different User

If the user doesn't have permissions:

```bash
scp pingxeno-agent-linux-amd64 root@192.168.1.100:/tmp/
```

## Quick Fix Checklist

1. ✅ Can you SSH to the server? `ssh cyb4x@192.168.1.100`
2. ✅ Is SSH service running on the server? `sudo systemctl status ssh`
3. ✅ Is port 22 open in firewall? `sudo ufw status`
4. ✅ Can you ping the server? `ping 192.168.1.100`
5. ✅ Do you have write permissions? Check `/tmp` directory

## Recommended Solution

If SSH works, try using `rsync` instead of `scp`:

```bash
rsync -avz --progress pingxeno-agent-linux-amd64 cyb4x@192.168.1.100:/tmp/pingxeno-agent
```

This is more reliable and shows progress.

