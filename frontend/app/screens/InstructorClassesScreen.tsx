import { FC } from "react"
import { Pressable, View, ViewStyle } from "react-native"
import { Screen } from "@/components/Screen"
import { Text } from "@/components/Text"
import { InstructorTabsNavigatorParamList } from "@/navigators/InstructorTabsNavigator"
import { NativeStackScreenProps } from "@react-navigation/native-stack"
import { useNavigation } from "@react-navigation/native"
import { Avatar, AvatarFallbackText, AvatarImage } from "@/components/ui/avatar"
import { Menu, MenuItem, MenuItemLabel } from "@/components/ui/menu"
import { Icon } from "@/components/ui/icon"
import { useAuth } from "@/context/AuthContext"
import { AppNavigation } from "@/navigators/AppNavigator"
import { SettingsIcon, CalendarCog, LogOut } from "lucide-react-native"
import { useSafeAreaInsets } from "react-native-safe-area-context"
import { Card } from "@/components/ui/card"

interface InstructorClassesScreenProps extends NativeStackScreenProps<InstructorTabsNavigatorParamList, "Classes"> { }

export const InstructorClassesScreen: FC<InstructorClassesScreenProps> = () => {
  // Pull in navigation via hook
  // const navigation = useNavigation<AppNavigation>()
  const navigation = useNavigation<AppNavigation>()
  const { profile, role, logout } = useAuth()
  const insets = useSafeAreaInsets()

  return (
    <View style={{ paddingTop: insets.top + 3 }}>
      <Menu
        placement="bottom left"
        offset={3}

        style={{ paddingTop: insets.top + 3 }}
        disabledKeys={["Settings"]}
        trigger={({ ...triggerProps }) => {
          return (
            <Pressable
              className="px-3"
              {...triggerProps}
            >
              <Avatar>
                <AvatarFallbackText>{profile?.name}</AvatarFallbackText>
                <AvatarImage source={{ uri: profile?.avatarUrl }} />
              </Avatar>
            </Pressable>
          )
        }}
      >

        <MenuItem key="Settings" textValue="Settings">
          <Icon as={SettingsIcon} size="sm" className="mr-2" />
          <MenuItemLabel size="sm">Settings</MenuItemLabel>
        </MenuItem>
        {
          role === "instructor" && <MenuItem key="student" textValue="هنرجو" onPress={() => navigation.navigate("Home", undefined)}>
            <Icon as={CalendarCog} size="sm" className="mr-2" />
            <MenuItemLabel size="sm">هنرجو</MenuItemLabel>
          </MenuItem>
        }
        <MenuItem key="logout" textValue="خروج" onPress={logout}>
          <Icon as={LogOut} size="sm" className="mr-2" />
          <MenuItemLabel size="sm">خروج</MenuItemLabel>
        </MenuItem>
      </Menu>
      <Card>

      </Card>
    </View >
  )
}

