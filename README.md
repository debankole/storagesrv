Assumptions:<p>
<l>
<li>Solution is using AWS FIFO queue and the data structure on server preserves the order of items, but items still may be processed out of order because of multiple parallel command processors, the assumption is that this problem is out of the scope</li>
<li>Server writes command results to the local file, there is no way for a client to get values</li>
<li>I assumed there is no need to cover everything with unit tests, so I covered only the ordered map as a sample</li>
<li>Binaries are compiled using "GOOS=linux GOARCH=arm64" env variables</li>
<l>
