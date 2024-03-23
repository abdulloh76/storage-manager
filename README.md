# storage-manager

### handle file upload
- split it to binary files
- store metada of file and where it is being sent to store
- send binary parts to storage-servers for store

### handle file download
- check db for file metadata
- get binary file parts from storage-servers
- combine them to whole file
- send it via http
