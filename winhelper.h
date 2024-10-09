#ifndef WINHELPER_H
#define WINHELPER_H

#include <combaseapi.h>

#ifdef __cplusplus
extern "C" {
#endif

	void freeUTF16String(LPWSTR str);

#ifdef __cplusplus
}
#endif

#endif // WINHELPER_H