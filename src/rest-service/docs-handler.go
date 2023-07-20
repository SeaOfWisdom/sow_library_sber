package rest

import (
	"fmt"
	"io"
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
		responError(w, http.StatusBadRequest, err.Error())

		return
	}
	defer file.Close()
	rs.logger.Infof("HandlerUploadDoc: uploaded File: %+v\n", handler.Filename)
	rs.logger.Infof("HandlerUploadDoc: file Size: %+v\n", handler.Size)
	rs.logger.Infof("HandlerUploadDoc: MIME Header: %+v\n", handler.Header)

	// read all of the contents of our uploaded file into a
	// byte array

	/// !!!! TODO !!!!
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		rs.logger.Errorf("HandlerUploadDoc: %v", err)
		responJSON(w, http.StatusOK, SuccessMsg{Msg: "OK"})

		return
	}

	imageResp, err := rs.ocrSrv.ExtractText(r.Context(), &ocr.ExtractTextRequest{
		Image: imageBytes,
	})
	if err != nil {
		rs.logger.Errorf("HandlerUploadDoc: while extract text via ocr service, err: %v", err)
		responJSON(w, http.StatusOK, SuccessMsg{Msg: "OK"})

		return
	}

	rs.logger.Errorf("HandlerUploadDoc: abstract: %v", imageResp.Abstract)
	rs.logger.Errorf("HandlerUploadDoc: main: %s", imageResp.Main)

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
