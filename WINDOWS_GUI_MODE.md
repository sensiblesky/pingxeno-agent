# Windows GUI Mode and Logging Features

## Overview
The PingXeno Agent now supports:
1. **Automatic Log Rotation** - Logs rotate when they reach 10MB
2. **GUI Mode** - Run agent in background mode (Windows)
3. **Background Execution** - Agent runs as a Windows service

## Log Rotation
- Log files automatically rotate when they reach **10MB**
- Old log files are renamed with a timestamp (e.g., `agent.log.20251230-123045`)
- Default log location:
  - Windows: `C:\ProgramData\PingXeno\agent.log`
  - Linux/macOS: `/var/log/pingxeno-agent.log`

## GUI Mode (Windows)
Run the agent in GUI/background mode:
```powershell
.\pingxeno-agent.exe gui
```

This will:
- Run the agent in the background
- Log all output to the log file (no console output)
- Continue running until stopped

## Windows Service Installation
The installation script automatically:
1. Creates a Windows service using NSSM (if available) or `sc.exe`
2. Configures the service to start automatically
3. Sets up logging to `C:\ProgramData\PingXeno\agent.log`
4. Runs the agent in GUI mode

### Manual Service Creation
If automatic service creation fails, you can create it manually:

**Using sc.exe:**
```powershell
sc.exe create PingXenoAgent binPath= "C:\Program Files\PingXeno\pingxeno-agent.exe gui --config C:\Program Files\PingXeno\agent.yaml" start= auto DisplayName= "PingXeno Monitoring Agent"
sc.exe start PingXenoAgent
```

**Using Start-Process (Background Process):**
```powershell
Start-Process -FilePath "C:\Program Files\PingXeno\pingxeno-agent.exe" -ArgumentList "gui --config C:\Program Files\PingXeno\agent.yaml" -WindowStyle Hidden
```

## Viewing Logs
View the latest logs:
```powershell
Get-Content C:\ProgramData\PingXeno\agent.log -Tail 50 -Wait
```

## Troubleshooting
If the agent hangs when running from CLI:
1. Use GUI mode: `.\pingxeno-agent.exe gui`
2. Check the log file for errors
3. Ensure the config file is properly formatted
4. Verify API URL and keys are correct

