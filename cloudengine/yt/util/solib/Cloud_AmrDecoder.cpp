#include "Cloud_AmrDecoder.h"  
#include "dec_if.h"
#include <iostream>  
using namespace std;  
 

#ifdef __cplusplus
extern "C" {
#endif
/*
 * Class:     Cloud_AmrDecoder
 * Method:    AmrDecoderInit
 * Signature: ()I
 */
JNIEXPORT jint JNICALL Java_Cloud_AmrDecoder_AmrDecoderInit
  (JNIEnv *, jobject)
{
	return (jint)D_IF_init();
}



/*
 * Class:     Cloud_AmrDecoder
 * Method:    AmrDecoderProcess
 * Signature: (I[SI[S)I
 */
JNIEXPORT jint JNICALL Java_Cloud_AmrDecoder_AmrDecoderProcess
  (JNIEnv * env, jobject, jint st , jint mode , jbyteArray env_inData , jint length, jshortArray env_outData)
{
	UWord8  * in = (UWord8 *)env->GetByteArrayElements(env_inData,NULL);
	Word16  * out = (Word16  *)env->GetShortArrayElements(env_outData,NULL);

	D_IF_decode((void *)st, in, out, 0);

	env->ReleaseByteArrayElements(  env_inData, (jbyte*)in, NULL );
	env->ReleaseShortArrayElements(  env_outData, (jshort*)out, NULL );

	return 0;

}

/*
 * Class:     Cloud_AmrDecoder
 * Method:    AmrDecoderExit
 * Signature: ()I
 */
JNIEXPORT jint JNICALL Java_Cloud_AmrDecoder_AmrDecoderExit
  (JNIEnv *, jobject, jint st)
{
	D_IF_exit((void *)st);
	return 0;
}


#ifdef __cplusplus
}
#endif
