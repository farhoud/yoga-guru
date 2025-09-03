import { createNativeStackNavigator } from "@react-navigation/native-stack"
import { LoginScreen } from "@/screens/AuthScreen/LoginScreen"
import { SignupScreen } from "@/screens/AuthScreen/SignupScreen"

export type AuthNavigatorParamList = {
  Login: undefined
  Signup: undefined
}

const Stack = createNativeStackNavigator<AuthNavigatorParamList>()
export const AuthNavigator = () => {
  return (
    <Stack.Navigator screenOptions={{ headerShown: false }}>
      <Stack.Screen name="Login" component={LoginScreen} />
      <Stack.Screen name="Signup" component={SignupScreen} />
    </Stack.Navigator>
  )
}
