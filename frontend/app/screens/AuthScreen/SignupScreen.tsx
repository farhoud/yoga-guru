import { FC, useEffect, useState } from "react"
import { View } from "react-native"
import type { AppNavigation, AppStackScreenProps } from "@/navigators/AppNavigator"
import { Text } from "@/components/Text"
import { Button, ButtonSpinner, ButtonText } from "@/components/ui/button"
import {
  FormControl,
  FormControlLabel,
  FormControlLabelText,
  FormControlError,
  FormControlErrorText,
} from "@/components/ui/form-control"
import { Heading } from "@/components/ui/heading"
import { Input, InputField } from "@/components/ui/input"
import { VStack } from "@/components/ui/vstack"
import { useNavigation } from "@react-navigation/native"
import { useAuth, SignupFormData, signupSchema } from "@/context/AuthContext"
import { zodResolver } from "@hookform/resolvers/zod"
import { Controller, useForm } from "react-hook-form"
import { NativeStackScreenProps } from "@react-navigation/native-stack"
import { AuthNavigatorParamList } from "@/navigators/AuthNavigator"
import { Select, SelectBackdrop, SelectContent, SelectDragIndicator, SelectDragIndicatorWrapper, SelectIcon, SelectInput, SelectItem, SelectPortal, SelectTrigger } from "@/components/ui/select"
import { ChevronDownIcon } from "@/components/ui/icon"
import { Toast, ToastDescription, useToast } from "@/components/ui/toast"

interface SignupScreenProps extends NativeStackScreenProps<AuthNavigatorParamList, "Signup"> { }

export const SignupScreen: FC<SignupScreenProps> = () => {
  // Pull in navigation via hook
  const navigation = useNavigation<AppNavigation>()

  const { signup, loading, error } = useAuth()

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<SignupFormData>({
    resolver: zodResolver(signupSchema),
    defaultValues: {
      name: "",
      phone: "+98",
      password: "",
    },
  })

  const toast = useToast()
  const [toastId, setToastId] = useState("0")

  const showNewToast = (error: string) => {
    console.log("show toest called")
    const newId = Math.random().toString()
    setToastId(newId)
    toast.show({
      id: newId,
      placement: "top",
      duration: 3000,
      render: ({ id }) => {
        const uniqueToastId = "toast-" + id
        return (
          <Toast nativeID={uniqueToastId} action="error" variant="solid">
            <ToastDescription>
              {error}
            </ToastDescription>
          </Toast>
        )
      },
    })
  }

  useEffect(() => {
    if (!!error) {
      showNewToast(error)
    }
  }, [error])


  return (
    <View className="bg-backgroundLight0 flex-1 justify-center px-4">
      <VStack space="lg" className="w-full max-w-sm self-center">
        <Heading size="lg" className="text-center">
          ثبت نام
        </Heading>
        {/* Phone number field */}
        <FormControl isInvalid={!!errors.name}>
          <FormControlLabel>
            <FormControlLabelText>نام</FormControlLabelText>
          </FormControlLabel>
          <Controller
            control={control}
            name="name"
            render={({ field: { onChange, onBlur, value } }) => (
              <Input variant="outline">
                <InputField
                  placeholder="نام و نام خانوادگی"
                  onBlur={onBlur}
                  onChangeText={onChange}
                  value={value}
                  keyboardType="phone-pad"
                />
              </Input>
            )}
          />
          {errors.name?.message && (
            <FormControlError>
              <FormControlErrorText>{errors.name.message}</FormControlErrorText>
            </FormControlError>
          )}
        </FormControl>

        {/* Phone number field */}
        <FormControl isInvalid={!!errors.phone}>
          <FormControlLabel>
            <FormControlLabelText>شماره تلفن</FormControlLabelText>
          </FormControlLabel>
          <Controller
            control={control}
            name="phone"
            render={({ field: { onChange, onBlur, value } }) => (
              <Input variant="outline">
                <InputField
                  placeholder="مثال: +989123456789"
                  onBlur={onBlur}
                  onChangeText={onChange}
                  value={value}
                  keyboardType="phone-pad"
                />
              </Input>
            )}
          />
          {errors.phone?.message && (
            <FormControlError>
              <FormControlErrorText>{errors.phone.message}</FormControlErrorText>
            </FormControlError>
          )}
        </FormControl>

        {/* Gender field */}
        <FormControl isInvalid={!!errors.gender} isRequired>
          <FormControlLabel>
            <FormControlLabelText>جنسیت</FormControlLabelText>
          </FormControlLabel>
          <Controller
            control={control}
            name="gender"
            render={({ field: { onChange, onBlur, value } }) =>
            (
              <Select onValueChange={onChange} selectedValue={value}>
                <SelectTrigger variant="outline" >
                  <SelectInput
                    onBlur={onBlur}
                    variant="outline"
                    className="p-2 w-[93%] rtl"
                    placeholder="جنسیت"
                  />
                  <SelectIcon as={ChevronDownIcon} />
                </SelectTrigger>
                <SelectPortal>
                  <SelectBackdrop />
                  <SelectContent>
                    <SelectDragIndicatorWrapper>
                      <SelectDragIndicator />
                    </SelectDragIndicatorWrapper>
                    <SelectItem label="مرد" value="male" />
                    <SelectItem label="زن" value="female" />
                  </SelectContent>
                </SelectPortal>
              </Select>
            )
            }
          />

          {errors.gender?.message && (
            <FormControlError>
              <FormControlErrorText>{errors.gender.message}</FormControlErrorText>
            </FormControlError>
          )}
        </FormControl>

        {/* Password field */}
        <FormControl isInvalid={!!errors.password}>
          <FormControlLabel>
            <FormControlLabelText>رمز عبور</FormControlLabelText>
          </FormControlLabel>
          <Controller
            control={control}
            name="password"
            render={({ field: { onChange, onBlur, value } }) => (
              <Input variant="outline">
                <InputField
                  placeholder="رمز عبور"
                  onBlur={onBlur}
                  onChangeText={onChange}
                  value={value}
                  secureTextEntry
                />
              </Input>
            )}
          />
          {errors.password?.message && (
            <FormControlError>
              <FormControlErrorText>{errors.password.message}</FormControlErrorText>
            </FormControlError>
          )}
        </FormControl>

        {/* Submit */}
        <Button onPress={handleSubmit(signup)} disabled={loading}>
          {loading && <ButtonSpinner />}
          <ButtonText>ثبت نام</ButtonText>
        </Button>

        <Text
          onPress={() => navigation.navigate("Auth", { screen: "Login" })}
          className="text-textLight400 mt-2 text-center"
        >
          قبلاً ثبت‌نام کرده‌اید؟ وارد شوید
        </Text>
      </VStack>
    </View>
  )
}
