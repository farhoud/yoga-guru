import { FC } from "react"
import { View, ViewStyle } from "react-native"
import type { AppStackScreenProps, AppNavigation } from "@/navigators/AppNavigator"
import { Screen } from "@/components/Screen"
import { Text } from "@/components/Text"
import { Box } from "@/components/ui/box"
import { ButtonText, Button } from "@/components/ui/button"
import { Heading } from "@/components/ui/heading"
import { HStack } from "@/components/ui/hstack"
import { VStack } from "@/components/ui/vstack"
import { Image } from "@/components/ui/image"
import { Card } from "@/components/ui/card"
import { useSafeAreaInsets } from "react-native-safe-area-context"
import { useNavigation } from "@react-navigation/native"
import { Avatar, AvatarFallbackText, AvatarImage } from "@/components/ui/avatar"
import { Menu, MenuItem, MenuItemLabel } from "@/components/ui/menu"
import { Icon, AddIcon, GlobeIcon, PlayIcon, SettingsIcon } from "@/components/ui/icon"
import { CalendarCog, LogOut } from "lucide-react-native"
import { Pressable } from "@/components/ui/pressable"
import { ImageBackground } from "@/components/ui/image-background"
import { useAuth } from "@/context/AuthContext"

interface HomeScreenProps extends AppStackScreenProps<"Home"> { }

const cardBackground = require("@assets/images/background.png")

export const HomeScreen: FC<HomeScreenProps> = () => {
  // Pull in navigation via hook
  const navigation = useNavigation<AppNavigation>()
  const { profile, role, logout } = useAuth()

  const insets = useSafeAreaInsets()

  const upcomingClasses = [
    {
      id: "1",
      name: "Morning Flow",
      time: "10:07 Ppm - 2:00 Ppm",
      icon: "🧘",
    },
    {
      id: "2",
      name: "Evening Yin",
      time: "12:07 Ppm - 2:30 Ppm",
      icon: "🧘",
    },
    {
      id: "3",
      name: "Power Vinyasa",
      time: "4:30 Ppm - 5:30 Ppm",
      icon: "🧘",
    },
  ]
  return (
    <Screen style={$root} preset="fixed">
      {/* Header Section */}
      <View style={{ minHeight: 500 }}>
        <ImageBackground
          imageClassName="rounded-b-3xl"
          style={{ flex: 1 }}
          resizeMode="cover"
          source={cardBackground}
          alt="sege"
        >
          <View className="flex-1 flex-col justify-between">
            <Menu
              placement="bottom left"
              offset={3}
              disabledKeys={["Settings"]}
              trigger={({ ...triggerProps }) => {
                return (
                  <Pressable
                    className="px-3"
                    style={{ paddingTop: insets.top + 3 }}
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
                role === "instructor" && <MenuItem key="instructor" textValue="برنامه" onPress={() => navigation.navigate("Instructor", { screen: "Classes" })}>
                  <Icon as={CalendarCog} size="sm" className="mr-2" />
                  <MenuItemLabel size="sm">برنامه</MenuItemLabel>
                </MenuItem>
              }
              <MenuItem key="logout" textValue="خروج" onPress={logout}>
                <Icon as={LogOut} size="sm" className="mr-2" />
                <MenuItemLabel size="sm">خروج</MenuItemLabel>
              </MenuItem>
            </Menu>
            {/* Quote Section */}

            {/* Main Action Buttons */}
            <Box className="p-5">
              <Heading size="md" className="text-center text-neutral-900">
                “Inhale the future, exhale the past”
              </Heading>
            </Box>
            <Box className="mt-auto px-10 py-2">
              <Button
                className="rounded-3xl bg-primary-600 align-bottom active:bg-primary-700"
                size="xl"
              >
                <ButtonText className="font-bold text-white">ورود به جلسه</ButtonText>
              </Button>
            </Box>
          </View>
        </ImageBackground>
      </View>

      {/* Upcoming Class Schedule */}
      <Box className="px-2 py-5">
        <Box className="rounded-3xl bg-white p-5">
          <HStack>
            <Heading size="md" className="mb-4 text-neutral-900">
              برنامه کلاس های آتی
            </Heading>
            <Box className="ml-auto">
              <Button
                className="rounded-3xl bg-secondary-600 align-bottom active:bg-secondary-700"
                size="sm"
                onPress={() => {
                  navigation.navigate("ExploreCourses")
                }}
              >
                <ButtonText className="font-bold text-white">لیست دوره ها</ButtonText>
              </Button>
            </Box>
          </HStack>
          {upcomingClasses.map((item) => (
            <HStack key={item.id} className="items-center justify-between border-neutral-200 py-2">
              <HStack className="items-center">
                <Box className="mr-3 h-10 w-10 items-center justify-center rounded-full bg-neutral-200">
                  <Text>{item.icon}</Text>
                </Box>
                <VStack>
                  <Text className="font-bold text-neutral-900">{item.name}</Text>
                  <Text className="text-neutral-500">{item.time}</Text>
                </VStack>
              </HStack>
              <Button className="rounded-3xl bg-primary-500 active:bg-primary-600">
                <ButtonText>ثبت نام</ButtonText>
              </Button>
            </HStack>
          ))}
        </Box>
      </Box>
    </Screen>
  )
}

const $root: ViewStyle = {
  flex: 1,
}
