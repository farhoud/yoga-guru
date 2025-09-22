import React, { FC, useState } from "react"
import { View } from "react-native"
import type { AppNavigation, AppStackScreenProps } from "@/navigators/AppNavigator"
import { VStack } from "@/components/ui/vstack"
import { ButtonText, Button, ButtonSpinner } from "@/components/ui/button"
import { Text } from "@/components/ui/text"
import {
  FormControl,
  FormControlLabel,
  FormControlLabelText,
  FormControlError,
  FormControlErrorText,
} from "@/components/ui/form-control"
import { Input, InputField } from "@/components/ui/input"
import { Heading } from "@/components/ui/heading"
import { useNavigation } from "@react-navigation/native"
import { useForm, Controller } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { LoginFormData, loginSchema, useAuth } from "@/context/AuthContext"
import { NativeStackScreenProps } from "@react-navigation/native-stack"
import { AuthNavigatorParamList } from "@/navigators/AuthNavigator"

interface LoginScreenProps extends NativeStackScreenProps<AuthNavigatorParamList, "Login"> { }

export const LoginScreen: FC<LoginScreenProps> = () => {
  // Pull in navigation via hook
  const navigation = useNavigation<AppNavigation>()

  const { login, loading } = useAuth()

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      phone: "+98",
      password: "",
    },
  })

  return (
    <View className="bg-backgroundLight0 flex-1 justify-center px-4">
      <VStack space="lg" className="w-full max-w-sm self-center">
        <Heading size="lg" className="text-center">
          ورود
        </Heading>

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
                  placeholder="مثال: 09123456789"
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

        {/* Submit button */}
        <Button disabled={loading} onPress={handleSubmit(login)}>
          {loading && <ButtonSpinner />}
          <ButtonText>ورود</ButtonText>
        </Button>

        <Text
          onPress={() => {
            navigation.navigate("Auth", { screen: "Signup" })
          }}
          className="text-textLight400 mt-2 text-center"
        >
          حساب کاربری ندارید؟ ثبت نام کنید
        </Text>
      </VStack>
    </View>
  )
}
