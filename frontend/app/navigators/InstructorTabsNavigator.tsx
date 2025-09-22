import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { InstructorActivitiesScreen } from "@/screens/InstructorActivitiesScreen"
import { InstructorHomeScreen } from "@/screens/InstructorHomeScreen"
import { Icon } from '@/components/ui/icon';
import { HomeIcon, Logs } from 'lucide-react-native';
import { Button, ButtonIcon } from '@/components/ui/button';

export type InstructorTabsNavigatorParamList = {
  Home: undefined
  Activities: undefined
}


export const InstructorTabsNavigator = createBottomTabNavigator<InstructorTabsNavigatorParamList>();

export function InstructorTabs() {
  return (
    <InstructorTabsNavigator.Navigator screenOptions={{
      headerShown: false, tabBarShowLabel: false,
      tabBarButton: (props) => {
        return <Button {...props} ref={undefined}>
          {props.children}
        </Button>
      }
    }}>
      <InstructorTabsNavigator.Screen
        name="Home"
        component={InstructorHomeScreen}
        options={{
          tabBarIcon: ({ focused, color, size }) => (
            <ButtonIcon as={HomeIcon} color={color} />
          ),
        }} />
      <InstructorTabsNavigator.Screen
        name="Activities"
        component={InstructorActivitiesScreen}
        options={{
          tabBarIcon: ({ focused, color, size }) => (
            <ButtonIcon as={Logs} color={color} />
          ),
        }}
      />
    </InstructorTabsNavigator.Navigator>
  )
}

