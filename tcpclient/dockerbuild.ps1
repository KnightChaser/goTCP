# Get the USERNAME from user input
$USERNAME = Read-Host "Enter USERNAME "
Write-Output "You're now @$USERNAME"

# Generate a 6-digit random hex string (to make every container have different container names)
$RandomHex = -join ((1..6) | ForEach-Object { Get-Random -Minimum 0 -Maximum 16 } | ForEach-Object { $_.ToString("X") })
$containerName = "gotcpclient${USERNAME}${RandomHex}"

# Build and run Docker container
docker build -t gotcpclient:0.1 --build-arg --no-cache .
docker run --rm --network goTCPnet --interactive --tty --name $containerName gotcpclient:0.1 ./tcpclient -username $USERNAME
