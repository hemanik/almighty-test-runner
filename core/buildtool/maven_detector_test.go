package buildtool_test

import (
    "os"
    "fmt"
    "path"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    . "github.com/almighty/almighty-test-runner/core/configuration"
    . "github.com/almighty/almighty-test-runner/core/buildtool"
)

var _ = Describe("Maven build tool", func() {

    Context("Tool detection", func() {

        It("should support maven when pom.xml is found in the root directory of the project", func() {
            path := "/tmp/build_folder"

            mavenDetector := Create(path, TestRunnerConfiguration{})
            defer createTmpDir(path).with("pom.xml").andDefer()(purge(path))

            Expect(mavenDetector.InUse()).To(BeTrue())
        })

        It("should not support maven when pom.xml not found in the root directory", func() {
            path := "/tmp/build_folder"
            mavenDetector := Create(path, TestRunnerConfiguration{})
            defer createTmpDir(path).with("build.xml").andDefer()(purge(path))

            Expect(mavenDetector.InUse()).To(BeFalse())
        })

    })

    Context("When maven is used", func() {

        var mavenDetector *Maven

        It("should have default build command", func() {
            mavenDetector = &Maven{}
            Expect(mavenDetector.DefaultCmd()).To(Equal("mvn clean install"))
        })

        It("should get build command from configuration", func() {
            mavenDetector = Create("",
                TestRunnerConfiguration{
                    BuildToolConfiguration: BuildToolConfiguration{Command: "mvn test"},
            })
            Expect(mavenDetector.Command).To(Equal("mvn test"))
        })

    })

})

// --- Test helpers

func purge(path string) func() {
    return (func() {
        os.RemoveAll(path)
    })
}

func createTmpDir(path string) file {
    if err := os.MkdirAll(path, 0755); err != nil {
        fmt.Errorf("Unable to create directory %s. Reason %s", path, err)
    }
    return file{path: path}
}

type fileCreator interface {
    andDefer() func(fn func())
    with(fileName string) fileCreator
}

type file struct {
    path string
}

func (f file) with(fileName string) fileCreator {
    file, _ := os.Create(path.Join(f.path, fileName))
    defer file.Close()
    return f
}

func (f file) andDefer() func(func()) {
    return func(f func()) {
        f()
    }
}
