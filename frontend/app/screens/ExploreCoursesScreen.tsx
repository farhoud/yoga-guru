import { FC, useEffect, useState } from "react"
import { View, ViewStyle } from "react-native"
import type { AppStackScreenProps, AppNavigation } from "@/navigators/AppNavigator"
import { Screen } from "@/components/Screen"
import { Text } from "@/components/Text"
import { ListView } from "@/components/ListView"
import { Course } from "@/services/api/types"
import { CourseListItem } from "@/components/CourseListItem"
import { api } from "@/services/api"
import { Box } from "@/components/ui/box"
import { HStack } from "@/components/ui/hstack"
import { Button, ButtonIcon } from "@/components/ui/button"
import { ArrowRightIcon } from "@/components/ui/icon"
import { useSafeAreaInsets } from "react-native-safe-area-context"
import { useNavigation } from "@react-navigation/native"

interface ExploreCoursesScreenProps extends AppStackScreenProps<"ExploreCourses"> {}

export const ExploreCoursesScreen: FC<ExploreCoursesScreenProps> = () => {
  // Pull in navigation via hook
  const navigation = useNavigation<AppNavigation>()
  const [courses, setCourses] = useState<Course[]>([])

  const inset = useSafeAreaInsets()

  function goBack() {
    if (navigation.canGoBack()) {
      navigation.goBack()
    }
    navigation.navigate("Home")
  }

  useEffect(() => {
    api.getCourses().then((items) => {
      setCourses(items)
    })
  }, [])

  return (
    <View style={{ paddingTop: inset.top }} className="h-full w-full bg-background-500">
      <ListView
        ListHeaderComponent={() => {
          return (
            <Box className="p-2">
              <HStack>
                <Button size="lg" className="rounded-full p-3.5" onPress={goBack}>
                  <ButtonIcon as={ArrowRightIcon} />
                </Button>
              </HStack>
            </Box>
          )
        }}
        stickyHeaderHiddenOnScroll={false}
        data={courses}
        renderItem={({ item }) => <CourseListItem {...item} />}
      />
    </View>
  )
}

const $root: ViewStyle = {
  flex: 1,
}
