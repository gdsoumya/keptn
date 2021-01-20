package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-test/deep"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	keptnutils "github.com/keptn/go-utils/pkg/lib/keptn"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/metadata"
)

func TestGetMetadataHandlerFunc(t *testing.T) {
	type args struct {
		params metadata.MetadataParams
		p      *models.Principal
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "Get metadata",
			args: args{
				params: metadata.MetadataParams{
					HTTPRequest: nil,
				},
				p: nil,
			},
			wantStatus: 200,
		},
	}

	_ = os.Setenv("SECRET_TOKEN", "testtesttesttesttest")

	returnedStatus := 200

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(returnedStatus)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	_ = os.Setenv("EVENTBROKER_URI", ts.URL)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMetadataHandlerFunc(tt.args.params, tt.args.p)

			verifyHTTPResponse(got, tt.wantStatus, t)
		})
	}
}

func Test_metadataHandler_getMetadata(t *testing.T) {

	clientSet := fake.NewSimpleClientset(
		getBridgeDeployment(),
	)
	type fields struct {
		k8sClient kubernetes.Interface
		logger    keptn.LoggerInterface
	}
	tests := []struct {
		name   string
		fields fields
		want   middleware.Responder
	}{
		{
			name: "get bridge deployment info from k8s",
			fields: fields{
				k8sClient: clientSet,
				logger:    keptnutils.NewLogger("", "", "api"),
			},
			want: &metadata.MetadataOK{
				Payload: &models.Metadata{
					Bridgeversion: "bridge:0.8.0",
					Keptnlabel:    "keptn",
					Keptnservices: nil,
					Keptnversion:  "develop",
					Namespace:     "keptn",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("POD_NAMESPACE", "keptn")
			tmpSwaggerFileName := "tmp-swagger.yaml"

			if err := ioutil.WriteFile(tmpSwaggerFileName, []byte(testSwaggerYaml), os.ModePerm); err != nil {
				fmt.Println(err.Error())
			}
			defer os.Remove(tmpSwaggerFileName)
			h := &metadataHandler{
				k8sClient:       tt.fields.k8sClient,
				logger:          tt.fields.logger,
				swaggerFilePath: tmpSwaggerFileName,
			}
			got := h.getMetadata()

			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("getMetadata() did not return expected value:")
				for _, d := range diff {
					t.Log(d)
				}
			}
		})
	}
}

const testSwaggerYaml = `---
swagger: "2.0"
info:
  title: keptn api
  version: develop`

func getBridgeDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "bridge",
			Namespace:   "keptn",
			Annotations: map[string]string{},
		},
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "bridge",
							Image: "keptn/bridge:0.8.0",
						},
					},
				},
			},
		},
	}
}
