import { FC } from "react"
import { ViewStyle } from "react-native"
import { Screen } from "@/components/Screen"
import { Text } from "@/components/Text"
import { InstructorTabsNavigatorParamList } from "@/navigators/InstructorTabsNavigator"
import { NativeStackScreenProps } from "@react-navigation/native-stack"
// import { useNavigation } from "@react-navigation/native"

interface InstructorHomeScreenProps extends NativeStackScreenProps<InstructorTabsNavigatorParamList, "Home"> { }

export const InstructorHomeScreen: FC<InstructorHomeScreenProps> = () => {
  // Pull in navigation via hook
  // const navigation = useNavigation<AppNavigation>()
  return (
    <Screen style={$root} preset="scroll">
      <Text text="instructorHome" />
    </Screen>
  )
}

const $root: ViewStyle = {
  flex: 1,
}
