import { FC } from "react"
import { ViewStyle } from "react-native"
import { Screen } from "@/components/Screen"
import { Text } from "@/components/Text"
import { InstructorTabsNavigatorParamList } from "@/navigators/InstructorTabsNavigator"
import { NativeStackScreenProps } from "@react-navigation/native-stack"
// import { useNavigation } from "@react-navigation/native"

interface InstructorActivitiesScreenProps extends NativeStackScreenProps<InstructorTabsNavigatorParamList, "Activities"> { }

export const InstructorActivitiesScreen: FC<InstructorActivitiesScreenProps> = () => {
  // Pull in navigation via hook
  // const navigation = useNavigation<AppNavigation>()
  return (
    <Screen style={$root} preset="scroll">
      <Text text="instructorActivities" />
    </Screen>
  )
}

const $root: ViewStyle = {
  flex: 1,
}
