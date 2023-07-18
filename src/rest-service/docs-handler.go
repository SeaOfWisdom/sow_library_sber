package rest

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	ocr "github.com/SeaOfWisdom/sow_proto/ocr-srv"
	"github.com/gorilla/mux"
)

// HandlerUploadDoc UploadDoc godoc
// @Summary      Upload doc of work
// @Description  Uploading documents confirming work
// @Tags         Docs
// @Accept       json
// @Produce      json
// @Param        doc_type   path      string  true  "work id"
// @Success      200  {object}  SuccessMsg
// @Failure      400  {object}  ErrorMsg
// @Security Bearer
// @Router       /upload_doc [put]
func (rs *RestSrv) HandlerUploadDoc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docType, ok := vars["doc_type"]
	if !ok {
		responError(w, http.StatusBadRequest, "null request param")
		return
	}
	rs.logger.Info(fmt.Sprintf("doc type: %s", docType))
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("doc")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		responError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// read all of the contents of our uploaded file into a
	// byte array

	/// !!!! TODO !!!!
	imageBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		responJSON(w, http.StatusOK, SuccessMsg{Msg: "OK"})
		return
	}

	imageResp, err := rs.ocrSrv.ExtractText(context.Background(), &ocr.ExtractTextRequest{
		Image: imageBytes,
	})
	if err != nil {
		rs.logger.Error(fmt.Sprintf("while extract text via ocr service, err: %v", err))
		responJSON(w, http.StatusOK, SuccessMsg{Msg: "OK"})
		return
	}

	fmt.Println("Abstract: ", imageResp.Abstract)
	fmt.Println("Main: ", imageResp.Main)

	// TODO
	if docType == "paper" {
		// if _, _, err := rs.libSrv.PublishWork(web3Address, request.Work); err != nil {
		// 	responError(w, http.StatusBadRequest, err.Error())
		// 	return
		// }
	}

	// // Create a temporary file within our temp-images directory that follows
	// // a particular naming pattern
	// tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer tempFile.Close()

	// // write this byte array to our temporary file
	// tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	responJSON(w, http.StatusOK, SuccessMsg{Msg: "OK"})
}
