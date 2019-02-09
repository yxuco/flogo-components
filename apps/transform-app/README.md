# transform-app
Flogo application for demonstrating the use of [jsmapper activity](https://github.com/yxuco/flogo-components/tree/master/activity/jsmapper)

## Build and start the app
The [Makefile](https://github.com/yxuco/flogo-components/tree/master/apps/transform-app/Makefile) contains all steps to build, run, and test the sample app.  Following command will rebuild and start the app:
```bash
make
```

## Test
It includes 2 tests for the 2 branches of the Flogo flow.  Run the tests by
```bash
make first
```
or
```bash
make list
```

## Import model to Flogo UI
To view or edit the application, you may import [mapper_app.json](https://github.com/yxuco/flogo-components/tree/master/apps/transform-app/mapper_app.json) into Flogo UI.