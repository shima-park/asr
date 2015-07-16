package util

//#cgo LDFLAGS: -ldl
//#cgo LDFLAGS: -lm
/*

#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>
#include "solib/typedef.h"

void myprint(char* s) {
        printf("%s", s);
}

//D_IF_init
void * (*amr_decode_init)( void);

//D_IF_decode
void * (*amr_decode)( void *st, UWord8 *bits, Word16 *synth, Word32 lfi);

//D_IF_exit
void * (*amr_decode_exit)(void *state);

//E_IF_encode
int  * (*amr_encode)(void *st, Word16 req_mode, Word16 *speech, UWord8 *serial, Word16 dtx);

//E_IF_init
void * (*amr_encode_init)(void);

//E_IF_exit
void * (*amr_encode_exit)(void *state);

void * libm_handle = NULL;



int CInit(){
     char *errorInfo;

     libm_handle = dlopen("libamr.so", RTLD_LAZY );
     if (!libm_handle){
         printf("Open Error:%s.\n",dlerror());
         return 0;
     }

     amr_encode_init = dlsym(libm_handle,"E_IF_init");
     errorInfo = dlerror();
     if (errorInfo != NULL){
         printf("Dlsym Error  :%s.\n",errorInfo);
         return 0;
     }

     amr_encode_exit = dlsym(libm_handle,"E_IF_exit");
     errorInfo = dlerror();
     if (errorInfo != NULL){
         printf("Dlsym Error  :%s.\n",errorInfo);
         return 0;
     }

     amr_encode = dlsym(libm_handle,"E_IF_encode");
     errorInfo = dlerror();
     if (errorInfo != NULL){
         printf("Dlsym Error  :%s.\n",errorInfo);
         return 0;
     }

     amr_decode_init = dlsym(libm_handle,"D_IF_init");
     errorInfo = dlerror();
     if (errorInfo != NULL){
         printf("Dlsym Error  :%s.\n",errorInfo);
         return 0;
     }

     amr_decode = dlsym(libm_handle,"D_IF_decode");
     errorInfo = dlerror();
     if (errorInfo != NULL){
         printf("Dlsym Error  :%s.\n",errorInfo);
         return 0;
     }

     amr_decode_exit = dlsym(libm_handle,"D_IF_exit");
     errorInfo = dlerror();
     if (errorInfo != NULL){
         printf("Dlsym Error  :%s.\n",errorInfo);
         return 0;
     }
     return 0;
 }

void CExit(){
     dlclose(libm_handle);
}

void *E_IF_init(void){
	return (*amr_encode_init)();
}

int *E_IF_encode(void *st, Word16 req_mode, Word16 *speech, UWord8 *serial, Word16 dtx){
	return (*amr_encode)(st,req_mode,speech,serial,dtx);
}

void E_IF_exit(void *state){
	(*amr_encode_exit)(state);
}

void *D_IF_init(void){
	return (*amr_decode_init)();
}

void D_IF_decode( void *st, UWord8 *bits, Word16 *synth, Word32 lfi){
	(*amr_decode)(st,bits,synth,lfi);
}

void D_IF_exit(void *state){
	(*amr_decode_exit)(state);
}
//


*/
import "C"

import (
	"unsafe"
)

func init() {
	C.CInit()
}

func Close() {
	C.CExit()
}

func DecoderFix(inputStream []byte) []byte {

	st := C.D_IF_init()

	defer C.D_IF_exit(st)

	var outputStream []byte

	read_len := 18

	//inputStream = inputStream[9:]

	for i := 0; i < (len(inputStream) / read_len); i++ {

		var serialBytes []byte

		var synthBytes []byte = make([]byte, 640)

		var start, end int

		if (i+1)*read_len > len(inputStream) {
			start = i * read_len
			end = i*read_len + (len(inputStream) - i*read_len)
		} else {
			start = i * read_len
			end = (i + 1) * read_len
		}

		serialBytes = inputStream[start:end]

		serial := (*C.UWord8)(unsafe.Pointer(&serialBytes[0]))

		synth := (*C.Word16)(unsafe.Pointer(&synthBytes[0]))

		C.D_IF_decode(st, serial, synth, 0)

		outputStream = append(outputStream, synthBytes...)
	}

	return outputStream
}
