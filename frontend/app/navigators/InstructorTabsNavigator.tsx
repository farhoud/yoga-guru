import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { InstructorActivitiesScreen } from "@/screens/InstructorActivitiesScreen"
import { InstructorClassesScreen } from "@/screens/InstructorClassesScreen"
import { Icon } from '@/components/ui/icon';
import { HomeIcon, Logs } from 'lucide-react-native';
import { Button, ButtonIcon } from '@/components/ui/button';

export type InstructorTabsNavigatorParamList = {
  Classes: undefined
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
        name="Classes"
        component={InstructorClassesScreen}
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

