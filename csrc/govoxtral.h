#ifndef GOVOXTRAL_H
#define GOVOXTRAL_H

#ifdef _WIN32
    #ifdef GOVOXTRAL_EXPORTS
        #define GOVOXTRAL_API __declspec(dllexport)
    #else
        #define GOVOXTRAL_API __declspec(dllimport)
    #endif
#else
    #define GOVOXTRAL_API
#endif

extern GOVOXTRAL_API int load_model(const char *model_dir);
extern GOVOXTRAL_API const char *transcribe(const char *wav_path);
extern GOVOXTRAL_API void free_result(void);

#endif /* GOVOXTRAL_H */
