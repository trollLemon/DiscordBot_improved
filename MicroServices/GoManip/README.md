# GoManip
CPU bound image manipulation service.

## Building

GoManip uses a wrapper around OpenCV, which requires building OpenCV from source. Building the Docker image with 
```bash
docker build -t GoManip .   
```

Will go through all build steps, including building OpenCV. 

The final image is build from a distroless image, and only includes the binary and the OpenCV (and other) libraries.


## Testing

To run tests build the testing image with
```bash
docker build --target tester -t GoManip-tests .  
```
and run the tests with
```bash
docker run  --rm GoManip-tests
```


You can also run the tests on your machine with 
```bash
go test -v ./...
```

However this may introduce issues if you have a different OpenCV version on your machine.

## API Endpoints
All endpoints expect to be called via POST, with parameters (if any) supplied via query params

The following list contains all supported endpoints for image manipulation functions and any query parameters:
- `/api/image/invert/`
- `/api/image/saturate/`
  - `saturation (float)`
- `/api/image/edgeDetection/`
  - `lower (int64)`
  - `higher (int64)`
- `/api/image/morphology/`
  - `kernelSize (int64)`
  - `iterations (int64)`
  - `type (string)`
- `/api/image/reduction/`
  - `quality (float)`
- `/api/image/text/`
  - `text (string)`
  - `fontScale (float)`
  - `xPerc (float)` (percentage along the x-axis of an image, i.e 0.5 for the middle along the width)
  - `xPerc (float)` (percentage along the y-axis of an image, i.e 0.5 for the middle along the height)
- `/api/image/randomFilter/`
  - `kernelSize (int64)`
  - `maxVal (int64)`
  - `minVal (int64)`
  - `normalize (bool)` whether or not to normalize the kernel
- `/api/image/shuffle/`
  - `partitions (int64)`


## Return Values
On successful operations, the api will return the result image as raw bytes in the HTTP body, with HTTP status code=200. For errors during processing,
i.e. invalid parameters, the api will return an error string json, with status code=400.



## Command Line Arguments
GoManip has two command line arguments:
 - `--pretty_print` to enable pretty printing rather than json in the logs. The default value is false.
 - `--num_workers`  to specify how many worker goroutines to spawn. The default value is the max number of logical cpus available to the process. 

Since all operations are vectorized due to opencv, image manipulation functions are fast, but can clog up the CPU if too many jobs are dispatched. `--num_workers` can help set a bound for how many
jobs will have threaded OpenCV operations.


